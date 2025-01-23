package internal

import (
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type EtcdClient struct {
	client *clientv3.Client
	ctx    context.Context
}

func NewEtcdClient(client *clientv3.Client, ctx context.Context) *EtcdClient {
	return &EtcdClient{
		client: client,
		ctx:    ctx,
	}
}

func (e *EtcdClient) Register(serviceName string, addr string) error {
	//TODO:chang to zap/logger
	log.Println("Registering service", serviceName)
	lease := clientv3.NewLease(e.client)
	_, cancel := context.WithTimeout(e.ctx, time.Second*5)
	defer cancel()

	_, err := lease.Grant(e.ctx, 30)
	if err != nil {
		return errors.New("Error in Grant :" + err.Error())
	}

	return nil
}
