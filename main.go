package main

import (
	"flag"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/pkg/grpc_service/interceptors"
	"github.com/peterouob/golang_template/server"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8082, "grpc service port")
)

func init() {
	tools.InitLogger()
	configs.InitViper()
}

func main() {
	flag.Parse()
	servers := []server.GrpcServer{
		server.RegisterUserService("echo", nil, nil),
		server.RegisterUserService("login", nil, nil),
		server.RegisterUserService("jwt", []grpc.UnaryServerInterceptor{interceptors.TokenInterceptors}, nil),
	}
	ports := []int{8081, 8082, 8083}
	for i, gserver := range servers {
		go gserver.InitServer(ports[i])
	}
	server.GrpcGatewayServer(*port)
}
