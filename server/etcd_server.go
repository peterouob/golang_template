package server

import (
	"context"
	"errors"
	"github.com/peterouob/golang_template/internal"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func InitEtcdServer(addr []string) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	})
	ctx := context.Background()

	c := internal.NewEtcdClient(client, ctx)
	if err != nil {
		panic(errors.New("create etcd client failed :" + err.Error()))
	}

	//TODO:need same with rpc port
	err = c.Register("go_", ":8081")
	if err != nil {
		panic(errors.New("register go_ failed :" + err.Error()))
	}
}
