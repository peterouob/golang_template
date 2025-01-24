package etcdservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/peterouob/golang_template/tools"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"sync"
	"time"
)

type EtcdService struct {
	client    *clientv3.Client
	heartbeat int64
}

var (
	serviceHub *EtcdService
	hubOnce    sync.Once
)

func RegisterETCD(etcdServers []string, heartbeat int64) *EtcdService {
	tools.Log("Starting etcd servers...")
	hubOnce.Do(func() {
		if serviceHub == nil {
			client, err := clientv3.New(clientv3.Config{
				Endpoints:   etcdServers,
				DialTimeout: 5 * time.Second,
			})
			tools.HandelError("new etcd client error", err)
			serviceHub = &EtcdService{client: client, heartbeat: heartbeat}
		}
	})

	return serviceHub
}

func (s *EtcdService) Register(service string, endpoint string, leaseID clientv3.LeaseID) clientv3.LeaseID {
	ctx := context.Background()
	if leaseID <= 0 {
		lease, err := s.client.Grant(ctx, s.heartbeat)
		tools.HandelError("grant lease error", err)
		key := fmt.Sprintf("%s/%s/%s",
			strings.TrimRight("/service/grpc", "/"),
			service,
			endpoint)

		_, err = s.client.Put(ctx, key, "", clientv3.WithLease(leaseID))
		tools.HandelError(fmt.Sprintf("puth in %s node %s on etcd error", service, endpoint), err)
		return lease.ID
	} else {
		_, err := s.client.KeepAliveOnce(ctx, leaseID)
		if errors.Is(err, rpctypes.ErrLeaseNotFound) {
			return s.Register(service, endpoint, leaseID)
		} else {
			tools.HandelError("keep lease error", err)
		}

		return leaseID
	}
}

func (s *EtcdService) UnRegister(service string, endpoint string) {
	ctx := context.Background()
	key := fmt.Sprintf("%s/%s/%s",
		strings.TrimRight("/service/grpc", "/"),
		service,
		endpoint)

	_, err := s.client.Delete(ctx, key)
	tools.HandelError("delete etcd node error", err)
	tools.Log(fmt.Sprintf("unregister %s node %s on etcd error", service, endpoint))
}
