package grpcclient

import (
	"context"
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/configs"
	etcdclient "github.com/peterouob/golang_template/pkg/etcd/client"
	grpcpool "github.com/peterouob/golang_template/pkg/grpc_service/pool"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"math"
	"math/rand"
	"sync"
	"time"
)

var (
	serverConn = sync.Map{}

	pool *grpcpool.Pool
)

func initPool(addr string, poolSize int) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt64), grpc.MaxCallSendMsgSize(math.MaxInt64)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}
	cfg := &configs.ClientConfig{}
	cfg.SetServerAddr(addr)
	cfg.SetPoolSize(poolSize)
	cfg.SetLifeTime(10 * time.Minute)
	cfg.SetLifeTimeDeviation(60 * time.Second)
	pool = grpcpool.NewPool(*cfg, opts, &grpcpool.RoundRobin{})
}

func GetGRPCClient(clientCfg *configs.EtcdGrpcCfg, serviceName string) (interface{}, *grpcpool.Pool, error) {
	if clientCfg == nil || tools.CheckStructNil(clientCfg) {
		clientCfg = &configs.EtcdGrpcCfg{}
		clientCfg.SetPoolSize(10)
		clientCfg.SetEndPoints([]string{"127.0.0.1:2379"})
		clientCfg.SetServiceName(serviceName)
	}

	hub := etcdclient.GetService(clientCfg.EndPoints)
	servers := hub.GetServiceEndPoint(clientCfg.ServiceName)
	if len(servers) == 0 {
		return nil, nil, fmt.Errorf("cannot get the service : %s", clientCfg.ServiceName)
	}

	server := servers[rand.Intn(len(servers))]
	tools.Log(fmt.Sprintf("from etcd connect to grpc server : %s", server))

	initPool(server, clientCfg.PoolSize)
	if client, exists := serverConn.Load(server); exists {
		return client, nil, nil
	}

	conn := pool.GetConn()

	_, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	switch serviceName {
	case "echo_service":
		return protobuf.NewEchoClient(conn), pool, nil
	default:
		return nil, nil, fmt.Errorf("unknown gRPC service: %s", serviceName)
	}
}
