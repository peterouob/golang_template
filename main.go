package main

import (
	"flag"
	"github.com/peterouob/golang_template/server"
	"github.com/peterouob/golang_template/tools"
)

var (
	servicePort = flag.Int("port", 8080, "grpc service port")
)

func init() {
	tools.InitLogger()
}

func main() {
	flag.Parse()

	server.InitGrpcServer(*servicePort)
}
