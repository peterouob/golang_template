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
	"google.golang.org/grpc/keepalive"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

var (
	serverConn = sync.Map{}
	mu         sync.Mutex
	pool       *grpc_service.Pool
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
	pool = grpc_service.NewPool(*cfg, opts, &grpc_service.RoundRobin{})
}

func EchoClient(clientCfg *configs.EtcdGrpcCfg) (protobuf.EchoClient, *grpc_service.Pool, error) {
	check := tools.CheckStructNil(clientCfg)
	if !check {
		clientCfg = &configs.EtcdGrpcCfg{}
		clientCfg.SetPoolSize(10)
		clientCfg.SetEndPoints([]string{"127.0.0.1:2379"})
		clientCfg.SetServiceName("echo_service")
	}

	hub := etcdclient.GetService(clientCfg.EndPoints)
	servers := hub.GetServiceEndPoint(clientCfg.ServiceName)
	if len(servers) == 0 {
		return nil, nil, fmt.Errorf("get service %s fail", clientCfg.ServiceName)
	}

	idx := rand.Intn(len(servers))
	server := servers[idx]
	tools.Log(fmt.Sprintf("Connecting to gRPC server from etcd: %s", server))
	initPool(server, clientCfg.PoolSize)
	log.Printf("%+v", pool)
	cc, exists := serverConn.Load(server)
	if !exists {
		mu.Lock()
		defer mu.Unlock()
		conn := pool.GetConn()
		client := protobuf.NewEchoClient(conn)
		serverConn.Store(server, client)
		return client, pool, nil
	}

	if client, ok := cc.(protobuf.EchoClient); ok {
		return client, nil, nil
	}

	return nil, nil, fmt.Errorf("client fail")
}
