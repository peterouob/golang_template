package server

import (
	"errors"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/service"
	"google.golang.org/grpc"
	"net"
)

func InitGrpcServer() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(errors.New("error listening port"))
	}
	s := grpc.NewServer()
	echo := service.NewEchoServer()
	protobuf.RegisterEchoServer(s, echo)

	if err := s.Serve(lis); err != nil {
		panic(errors.New("error serving"))
	}
}
