package etcdregister

import (
	"fmt"
	etcdservice "github.com/peterouob/golang_template/pkg/etcd/server"
	"github.com/peterouob/golang_template/tools"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type EtcdRegister struct {
	client  *etcdservice.EtcdService
	leaseId clientv3.LeaseID
	heart   int64
}

func NewEtcdRegister(endpoints []string, heart int64) *EtcdRegister {
	return &EtcdRegister{
		client: etcdservice.RegisterETCD(endpoints, heart),
		heart:  heart,
	}
}

func (e *EtcdRegister) Register(serviceName, addr string) {
	e.leaseId = e.client.Register(serviceName, addr, 0)
	tools.Log(fmt.Sprintf("Registered service %s at %s", serviceName, addr))
	go func() {
		for {
			e.client.Register(serviceName, addr, e.leaseId)
			time.Sleep(time.Duration(e.heart)*time.Second - 100*time.Millisecond)
		}
	}()
	e.listenExit(serviceName, addr)
}

func (e *EtcdRegister) UnRegister(serviceName, addr string) {
	e.client.UnRegister(serviceName, addr)
	tools.Log(fmt.Sprintf("unregiter service: %s from etcd, addr: %s", serviceName, addr))
}

func (e *EtcdRegister) listenExit(serviceName, addr string) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		tools.Log(fmt.Sprintf("receive a signal %s", sig.String()))
		e.client.UnRegister(serviceName, addr)
		os.Exit(0)
	}()
}
