package user

import (
	"context"
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenTestServer struct {
	protobuf.UnimplementedUserServer
}

func NewTokenTestServer() *TokenTestServer {
	return &TokenTestServer{}
}

// TokenTest I think maybe this can abstract to middleware and be a service
func (t TokenTestServer) TokenTest(ctx context.Context, in *protobuf.TokenTestRequest) (*protobuf.TokenTestResponse, error) {
	userId, ok := ctx.Value("uid").(int64)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "userId not found in context")
	}
	return &protobuf.TokenTestResponse{
		Msg: fmt.Sprintf("This is Token Test! your id is :%d", userId),
	}, nil
}
