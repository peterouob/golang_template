package main

import (
	"flag"
	"github.com/peterouob/golang_template/server"
	"github.com/peterouob/golang_template/tools"
)

var (
	servicePort = flag.Int("port", 8080, "grpc service port")
)

func main() {
	flag.Parse()
	tools.InitLogger()
	server.InitGrpcServer(*servicePort)
}
