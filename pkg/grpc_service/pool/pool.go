package grpcpool

import (
	"errors"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"log"
)

type Pool struct {
	Cfg         configs.ClientConfig
	Conns       []*PoolConn
	DialOpts    []grpc.DialOption
	LoadBalance LoadBalance
}

func NewPool(cfg configs.ClientConfig, dialOpts []grpc.DialOption, loadBalance LoadBalance) *Pool {
	if cfg.PoolSize <= 0 {
		tools.HandelError("error in NewPool", errors.New("the pool size is smaller than zero"))
	}
	if len(cfg.ServerAddr) == 0 {
		tools.HandelError("error in NewPool", errors.New("the server address is empty"))
	}

	conns := make([]*PoolConn, 0, cfg.PoolSize)
	for i := 0; i < cfg.PoolSize; i++ {
		conn := new(PoolConn)
		conn.Refresh(cfg, dialOpts...)
		conns = append(conns, conn)
	}

	return &Pool{
		Cfg:         cfg,
		Conns:       conns,
		DialOpts:    dialOpts,
		LoadBalance: loadBalance,
	}
}

func (p *Pool) GetConn() *grpc.ClientConn {
	idx := p.LoadBalance.Select(len(p.Conns))
	conn := p.Conns[idx]

	log.Printf("%+v", conn)

	if conn.ShouldRefresh() {
		conn.Refresh(p.Cfg, p.DialOpts...)
	}

	if conn != nil && isConnectionHealthy(conn.conn) {
		return conn.conn
	} else {
		return p.getNextConn(idx, len(p.Conns))
	}

}

func (p *Pool) getNextConn(currIdx, curSize int) *grpc.ClientConn {
	i := currIdx
	for {
		i = (i + 1) % curSize
		if i == currIdx {
			tools.Log("no available connection ... reconnect on ")
			return p.connect()
		}
		if p.Conns[i] != nil && isConnectionHealthy(p.Conns[i].conn) {
			tools.Log("find connection in pool")
			return p.Conns[i].conn
		}
	}
}

func (p *Pool) connect() *grpc.ClientConn {
	conn, err := grpc.NewClient(p.Cfg.ServerAddr, p.DialOpts...)
	tools.HandelError("error in connect at pool", err)
	return conn
}
