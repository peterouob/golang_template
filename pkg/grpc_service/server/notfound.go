package grpcserver

import (
	"context"
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"google.golang.org/grpc"
)

type NotFoundServer struct {
	protobuf.UnimplementedUserServer
}

func NewNotFoundServer() *NotFoundServer {
	return &NotFoundServer{}
}

func NotFound(ctx context.Context, in *protobuf.NotFoundRequest, opts ...grpc.CallOption) (*protobuf.NotFoundResponse, error) {
	return &protobuf.NotFoundResponse{
		Msg: fmt.Sprintf("Not found server %v", in.String()),
	}, nil
}
