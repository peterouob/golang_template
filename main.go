package main

import (
	"flag"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/server"
	"github.com/peterouob/golang_template/tools"
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
		server.RegisterEchoServer(),
		server.RegisterLoginServer(),
		server.RegisterTokenTestServer(),
	}
	ports := []int{8081, 8082, 8083}
	for i, gserver := range servers {
		go gserver.InitServer(ports[i])
	}
	//go server.InitLoginServer()
	//go server.InitTokenServer()

	server.GrpcGatewayServer(*port)
}
