package service

import (
	"context"
	"github.com/peterouob/golang_template/api/protobuf"
)

type EchoServer struct {
	protobuf.UnimplementedEchoServer
}

func NewEchoServer() *EchoServer {
	return &EchoServer{}
}

func (s EchoServer) Echo(ctx context.Context, in *protobuf.EchoRequest) (*protobuf.EchoResponse, error) {
	return &protobuf.EchoResponse{Name: "Server say hello to " + in.GetName()}, nil
}
