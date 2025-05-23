package etcdservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/peterouob/golang_template/utils"
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

// RegisterETCD TODO:服務連線失敗降級並等待etcd重新註冊上
func RegisterETCD(etcdServers []string, heartbeat int64) *EtcdService {
	hubOnce.Do(func() {
	start:
		if serviceHub == nil {
			client, err := clientv3.New(clientv3.Config{
				Endpoints:   etcdServers,
				DialTimeout: 5 * time.Second,
			})

			if err != nil {
				time.Sleep(5 * time.Second)
				utils.Log("wait for etcd servers to be ready...")
				goto start
			}

			serviceHub = &EtcdService{client: client, heartbeat: heartbeat}
		} else {
			serviceHub = &EtcdService{
				client:    serviceHub.client,
				heartbeat: heartbeat,
			}
		}
	})

	return serviceHub
}

func (s *EtcdService) Register(service string, endpoint string, leaseID clientv3.LeaseID) clientv3.LeaseID {
	ctx := context.Background()
	if leaseID <= 0 {
		lease, err := s.client.Grant(ctx, s.heartbeat)
		utils.HandelError("grant lease error", err)
		key := fmt.Sprintf("%s/%s/%s",
			strings.TrimRight("/service/grpc", "/"),
			service,
			endpoint)

		_, err = s.client.Put(ctx, key, "", clientv3.WithLease(leaseID))
		utils.HandelError(fmt.Sprintf("puth in %s node %s on etcd error", service, endpoint), err)
		return lease.ID
	} else {
		_, err := s.client.KeepAliveOnce(ctx, leaseID)
		if errors.Is(err, rpctypes.ErrLeaseNotFound) {
			return s.Register(service, endpoint, leaseID)
		} else {
			utils.HandelError("keep lease error", err)
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
	resp, err := s.client.Get(ctx, key)
	if err != nil || len(resp.Kvs) == 0 {
		utils.Log(fmt.Sprintf("Key %s not found in etcd", key))
		return
	}

	leaseID := clientv3.LeaseID(resp.Kvs[0].Lease)
	utils.Log(fmt.Sprintf("Revoking lease %d for key %s", leaseID, key))

	_, err = s.client.Revoke(ctx, leaseID)
	utils.HandelError("revoke lease error", err)

	_, err = s.client.Delete(ctx, key)
	utils.HandelError("delete etcd node error", err)
	utils.Log(fmt.Sprintf("unregistered %s node %s from etcd", service, endpoint))
}
