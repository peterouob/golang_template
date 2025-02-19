package grpcclient

import (
	"context"
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/configs"
	etcdclient "github.com/peterouob/golang_template/pkg/etcd/client"
	grpcpool "github.com/peterouob/golang_template/pkg/grpc_service/pool"
	"github.com/peterouob/golang_template/tools"
	"math/rand"
	"sync"
	"time"
)

var (
	serverConn = sync.Map{}
)

func GetGRPCClient(clientCfg *configs.EtcdGrpcCfg, serviceName string) (interface{}, error) {
	hub := etcdclient.GetService(clientCfg.EndPoints)
	servers := hub.GetServiceEndPoint(clientCfg.ServiceName)
	if len(servers) == 0 {
		return nil, fmt.Errorf("cannot get the service : %s", clientCfg.ServiceName)
	}

	server := servers[rand.Intn(len(servers))]
	tools.Log(fmt.Sprintf("from etcd connect to grpc server : %s", server))
	pool := grpcpool.InitPool(server, clientCfg.PoolSize)
	if client, exists := serverConn.Load(server); exists {
		return client, nil
	}

	conn := pool.GetConn()

	_, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	switch serviceName {
	case "echo_service":
		return protobuf.NewEchoClient(conn), nil
	case "login_service":
		return protobuf.NewUserClient(conn), nil
	case "tokentest_service":
		return protobuf.NewUserClient(conn), nil
	default:
		return nil, fmt.Errorf("unknown gRPC service: %s", serviceName)
	}
}
