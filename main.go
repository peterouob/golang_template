package main

import (
	"flag"
	"github.com/peterouob/golang_template/server"
	"github.com/peterouob/golang_template/tools"
)

var (
	port = flag.Int("port", 8080, "grpc service port")
)

func init() {
	tools.InitLogger()
}

func main() {
	flag.Parse()
	go func() {
		server.InitGrpcServer(*port)
	}()

	server.GrpcGatewayServer(*port)
}
