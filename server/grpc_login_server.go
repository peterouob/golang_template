package server

import (
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/grpc_service/interceptors"
	grpcserver "github.com/peterouob/golang_template/pkg/grpc_service/server"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
)

type LoginServer struct {
	BaseServer
}

func RegisterLoginServer() *EchoServer {
	echogrpc := &EchoServer{
		BaseServer{
			ServiceName: "login_service",
			RegisterFunc: func(server *grpc.Server) {
				login := grpcserver.NewLoginServer()
				protobuf.RegisterUserServer(server, login)
				tools.Log("register echo login success")
			},
		},
	}
	echogrpc.RegisterInterceptors(interceptors.LoggingInterceptor)
	return echogrpc
}
