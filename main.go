package main

import (
	"github.com/peterouob/golang_template/server"
	"github.com/peterouob/golang_template/tools"
)

func main() {
	tools.InitLogger()

	server.InitGrpcServer()

}
