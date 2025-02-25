package grpcserver

import (
	"context"
	"github.com/peterouob/golang_template/api/protobuf"
	mdb "github.com/peterouob/golang_template/pkg/db/mysql"
	"github.com/peterouob/golang_template/pkg/repository"
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
	user := &mdb.UserModel{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}
	u := repository.GetUserRepo()
	id, name := u.LoginUserByEmailAndPassword(*user)
	if id != -1 {
		aToken, rToken := verify.SaveToken(ctx, id)
		return &protobuf.LoginUserResponse{
			Name:         name,
			AccessToken:  aToken,
			RefreshToken: rToken,
		}, nil
	}
	return nil, status.Error(codes.Internal, "Error")
}
