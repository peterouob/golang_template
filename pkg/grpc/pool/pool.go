package pool

import (
	"errors"
	"fmt"
	"github.com/peterouob/golang_template/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Pool interface {
	Get() (Conn, error)
	Close() error
	Status() string
}

type pool struct {
	index     atomic.Uint32
	curr      atomic.Int32
	ref       atomic.Int32
	opt       configs.Option
	conns     []*conn
	addr      string
	closed    atomic.Int32
	checkTime time.Duration
	sync.RWMutex
}

func New(addr string, opt configs.Option) Pool {
	p := &pool{
		opt:       opt,
		conns:     make([]*conn, opt.MaxActive),
		addr:      addr,
		checkTime: time.Minute,
	}
	p.curr.Store(int32(opt.MaxActive))

	for i := 0; i < opt.MaxIdle; i++ {
		c, err := opt.Dial(addr)
		if err != nil {
			panic(fmt.Sprintf("error in dial %s , %v", addr, err.Error()))
		}
		p.conns[i] = p.wrapconn(c, false)
	}
	log.Printf("new pool success %v\n", p.Status())
	p.checkHealthy()
	return p
}

func (p *pool) Get() (Conn, error) {
	p.RLock()
	cur := p.curr.Load()
	p.RUnlock()

	p.incRef()
	nextRef := p.curr.Load()

	if cur == 0 {
		return nil, errors.New("pool closed")
	}

	if nextRef <= cur*int32(p.opt.MaxConcurrentStreams) {
		next := p.index.Add(1) % uint32(cur)
		return p.conns[next], nil
	}

	if cur == int32(p.opt.MaxActive) {
		if p.opt.Reuse {
			next := p.index.Add(1) % uint32(cur)
			return p.conns[next], nil
		}
		c, err := p.opt.Dial(p.addr)
		return p.wrapconn(c, true), err
	}

	p.Lock()
	cur = p.curr.Load()
	if cur < int32(p.opt.MaxActive) && nextRef > cur*int32(p.opt.MaxConcurrentStreams) {
		inc := cur
		if cur+inc > int32(p.opt.MaxActive) {
			inc = int32(p.opt.MaxActive) - cur
		}

		var i int32
		var err error
		for i = 0; i < inc; i++ {
			c, er := p.opt.Dial(p.addr)
			if er != nil {
				err = er
				break
			}
			p.reset(int(cur + i))
			p.conns[cur+i] = p.wrapconn(c, false)
		}

		cur += i
		log.Printf("grow pool: %d ---> %d, increment: %d, maxActive: %d\n",
			p.curr.Load(), cur, inc, p.opt.MaxActive)
		p.curr.Store(cur)
		if err != nil {
			p.Unlock()
			return nil, err
		}
	}
	p.Unlock()
	next := p.index.Add(1) % uint32(cur)
	return p.conns[next], nil
}
func (p *pool) Close() error {
	p.index.Store(0)
	p.curr.Store(0)
	p.ref.Store(0)
	p.closed.Store(1)
	p.delete(0)
	log.Printf("pool closed")
	return nil
}

func (p *pool) Status() string {
	return fmt.Sprintf("address:%s, index:%d, current:%d, ref:%d. option:%v",
		p.addr, p.index.Load(), p.curr.Load(), p.ref.Load(), p.opt)
}

func (p *pool) incRef() {
	p.ref.Add(1)
	if p.ref.Load() == math.MaxInt32 {
		panic(fmt.Sprint("ref overflow"))
	}
}

func (p *pool) decRef() {
	newRef := p.ref.Add(-1)
	if newRef < 0 && p.closed.Load() == 0 {
		panic(fmt.Sprint("ref overflow to negative"))
	}

	if newRef == 0 && p.curr.Load() > int32(p.opt.MaxIdle) {
		p.Lock()
		if p.ref.Load() == 0 {
			log.Printf("shrink pool: %d ---> %d, decrement: %d, maxActive: %d\n",
				p.curr.Load(), p.opt.MaxIdle, p.curr.Load()-int32(p.opt.MaxIdle), p.opt.MaxActive)

			p.curr.Store(int32(p.opt.MaxIdle))
			p.delete(p.opt.MaxIdle)
		}
		p.Unlock()
	}
}

func (p *pool) reset(idx int) {
	conn := p.conns[idx]
	if conn == nil {
		return
	}
	err := conn.reset()
	if err != nil {
		log.Printf("reset pool conn err:%v\n", err)
	}
	p.conns[idx] = nil
}

func (p *pool) delete(begin int) {
	for i := begin; i < p.opt.MaxActive; i++ {
		p.reset(i)
	}
}

func (p *pool) wrapconn(cc *grpc.ClientConn, once bool) *conn {
	return &conn{cc: cc, pool: p, once: once}
}

func (p *pool) checkHealthy() {
	go func() {
		for {
			if p.closed.Load() == 1 {
				return
			}

			select {
			case <-time.After(p.checkTime):
				p.reConnect()
			}
		}
	}()
}

func (p *pool) reConnect() {
	for i, conn := range p.conns {
		if conn == nil {
			continue
		}

		if conn.Value() == nil || conn.Value().GetState() == connectivity.Shutdown {
			log.Printf("reconnect pool conn[%d] is shutdown", i)
			p.Lock()
			_ = conn.reset()
			newConn, err := p.opt.Dial(p.addr)
			if err != nil {
				log.Printf("error in reconnect pool conn[%d]: %v", i, err)
			}
			p.conns[i] = p.wrapconn(newConn, false)
			p.Unlock()
		}
	}
}
