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
	//server.InitGrpcServer(*port)
	go server.InitLoginServer()
	go server.InitTokenServer()

	server.GrpcGatewayServer(*port)
}
