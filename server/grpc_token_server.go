package server

import (
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/grpc_service/interceptors"
	grpcserver "github.com/peterouob/golang_template/pkg/grpc_service/server"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
)

type TokenTestServer struct {
	BaseServer
}

func RegisterTokenTestServer() *EchoServer {
	echogrpc := &EchoServer{
		BaseServer{
			ServiceName: "tokentest_service",
			Registerfunc: func(server *grpc.Server) {
				tts := grpcserver.NewTokenTestServer()
				protobuf.RegisterUserServer(server, tts)
				tools.Log("register echo token test server success")
			},
		},
	}
	echogrpc.RegisterInterceptors(interceptors.LoggingInterceptor, interceptors.TokenInterceptors)
	return echogrpc
}
