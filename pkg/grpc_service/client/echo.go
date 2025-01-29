package grpcclient

import (
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/configs"
	etcdclient "github.com/peterouob/golang_template/pkg/etcd/client"
	"github.com/peterouob/golang_template/pkg/grpc_service"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math"
	"math/rand"
	"sync"
	"time"
)

var (
	serverConn = sync.Map{}
	mu         sync.Mutex
)

func initPool(addr string, poolSize int) *grpc_service.Pool {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt64), grpc.MaxCallSendMsgSize(math.MaxInt64)),
	}
	cfg := &configs.ClientConfig{}
	cfg.SetServerAddr(addr)
	cfg.SetPoolSize(poolSize)
	cfg.SetLifeTime(10 * time.Minute)
	cfg.SetLifeTimeDeviation(60 * time.Second)
	clientPool := grpc_service.NewPool(*cfg, opts, &grpc_service.RoundRobin{})
	return clientPool
}

// EchoClient TODO:Default Value Have Error
func EchoClient(clientCfg *configs.EtcdGrpcCfg) (protobuf.EchoClient, *grpc_service.Pool, error) {
	if clientCfg == nil {
		clientCfg = &configs.EtcdGrpcCfg{}
		clientCfg.SetPoolSize(10)
		clientCfg.SetEndPoints([]string{"127.0.0.1:2379"})
		clientCfg.SetServiceName("echo_service")
	}
	tools.Log(clientCfg.ServiceName)
	hub := etcdclient.GetService(clientCfg.EndPoints)
	servers := hub.GetServiceEndPoint(clientCfg.ServiceName)
	if len(servers) == 0 {
		return nil, nil, fmt.Errorf("get service %s fail", clientCfg.ServiceName)
	}

	idx := rand.Intn(len(servers))
	server := servers[idx]
	tools.Log(fmt.Sprintf("Connecting to gRPC server from etcd: %s", server))

	cc, exists := serverConn.Load(server)
	if !exists {
		mu.Lock()
		defer mu.Unlock()
		conn, err := grpc.NewClient(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
		tools.HandelError("new client fail", err)
		client := protobuf.NewEchoClient(conn)
		serverConn.Store(server, client)
		pool := initPool(server, clientCfg.PoolSize)
		return client, pool, nil
	}

	if client, ok := cc.(protobuf.EchoClient); ok {
		return client, nil, nil
	}

	return nil, nil, fmt.Errorf("client fail")
}
