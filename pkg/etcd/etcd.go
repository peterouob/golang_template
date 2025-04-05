package etcdregister

import (
	"fmt"
	etcdservice "github.com/peterouob/golang_template/pkg/etcd/server"
	"github.com/peterouob/golang_template/tools"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type EtcdRegister struct {
	client  *etcdservice.EtcdService
	leaseId clientv3.LeaseID
	heart   int64
}

func NewEtcdRegister(endpoints []string, heart int64) *EtcdRegister {
	c := etcdservice.RegisterETCD(endpoints, heart)
	e := &EtcdRegister{
		client: c,
		heart:  heart,
	}
	return e
}

func (e *EtcdRegister) Register(serviceName, addr string) {
	e.leaseId = e.client.Register(serviceName, addr, 0)
	//tools.Log(fmt.Sprintf("Registered service %s at %s", serviceName, addr))
	go func() {
		for {
			e.client.Register(serviceName, addr, e.leaseId)
			time.Sleep(time.Duration(e.heart)*time.Second - 100*time.Millisecond)
		}
	}()
}

func (e *EtcdRegister) UnRegister(serviceName, addr string) {
	e.client.UnRegister(serviceName, addr)
	tools.Log(fmt.Sprintf("unregiter service: %s from etcd, addr: %s", serviceName, addr))
}
