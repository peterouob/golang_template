package server

import (
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/grpc_service/interceptors"
	grpcserver "github.com/peterouob/golang_template/pkg/grpc_service/server"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
)

type EchoServer struct {
	BaseServer
}

func RegisterEchoServer() *EchoServer {
	echogrpc := &EchoServer{
		BaseServer{
			ServiceName: "echo_service",
			Registerfunc: func(server *grpc.Server) {
				echo := grpcserver.NewEchoServer()
				protobuf.RegisterEchoServer(server, echo)
				tools.Log("register echo server success")
			},
		},
	}
	echogrpc.RegisterInterceptors(interceptors.LoggingInterceptor)
	return echogrpc
}
