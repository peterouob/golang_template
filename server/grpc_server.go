package server

import (
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/service"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"net"
)

func InitGrpcServer() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	echo := service.NewEchoServer()
	protobuf.RegisterEchoServer(s, echo)

	tools.LogMessage("start grpc server ... \n")

	err = s.Serve(lis)
	if err != nil {
		tools.LogError("grpc server start error", err)
	}
}
