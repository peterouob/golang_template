package user

import (
	"context"
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
)

type NotFoundServer struct {
	protobuf.UnimplementedNotFoundServer
}

func NewNotFoundServer() *NotFoundServer {
	return &NotFoundServer{}
}

func (n NotFoundServer) NotFound(ctx context.Context, in *protobuf.NotFoundRequest) (*protobuf.NotFoundResponse, error) {
	return &protobuf.NotFoundResponse{
		Msg: fmt.Sprintf("Not found server %v", in.String()),
	}, nil
}
