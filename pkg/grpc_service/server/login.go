package grpcserver

import (
	"context"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/verify"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginServer struct {
	protobuf.UnimplementedUserServer
}

func NewLoginServer() *LoginServer {
	return &LoginServer{}
}

func (l LoginServer) LoginUser(ctx context.Context, in *protobuf.LoginUserRequest) (*protobuf.LoginUserResponse, error) {
	if in.Email != "admin" || in.Password != "admin" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument")
	}
	token := verify.NewToken(123)
	token.CreateToken()
	token.CreateRefreshToken()
	return &protobuf.LoginUserResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
