package grpcclient

import (
	"errors"
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	clientclient "github.com/peterouob/golang_template/pkg/etcd/client"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"sync"
)

var serverConn = sync.Map{}

func EchoClient() protobuf.EchoClient {
	hub := clientclient.GetService([]string{"127.0.0.1:2379"})
	servers := hub.GetServiceEndPoint("echo_service")
	if len(servers) == 0 {
		tools.HandelError("error in get client from etcd", errors.New("no server found"))
		return nil
	} else {
		idx := rand.Intn(len(servers))
		server := servers[idx]
		tools.Log(fmt.Sprintf("get client from etcd: %s", server))
		client, exists := serverConn.Load(server)
		if !exists {
			conn, err := grpc.NewClient(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
			tools.HandelError("error in new grpc client", err)
			c := protobuf.NewEchoClient(conn)
			serverConn.Store(server, c)
			return c
		}
		return client.(protobuf.EchoClient)
	}
}
