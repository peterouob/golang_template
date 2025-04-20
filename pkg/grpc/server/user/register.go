package user

import (
	"context"
	"github.com/peterouob/golang_template/api/protobuf"
	mdb "github.com/peterouob/golang_template/pkg/db/mysql"
	"github.com/peterouob/golang_template/pkg/repository"
	"github.com/peterouob/golang_template/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type RegisterServer struct {
	protobuf.UnimplementedUserServer
}

func NewRegisterServer() *RegisterServer {
	return &RegisterServer{}
}

func (r RegisterServer) RegisterUser(_ context.Context, in *protobuf.RegisterUserRequest) (*protobuf.RegisterUserResponse, error) {
	node := utils.NewIdWorker(10)
	uid := node.GenID()
	id, err := strconv.ParseInt(strconv.FormatInt(int64(uid), 10), 10, 64)
	utils.HandelError("convert string to int64 error", err)
	user := &mdb.UserModel{
		Name:     in.GetName(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
		Id:       id,
	}
	u := repository.GetUserRepo()
	if ok := u.CreateUser(*user); !ok {
		return &protobuf.RegisterUserResponse{}, status.Error(codes.Internal, "create user error")
	}
	return &protobuf.RegisterUserResponse{
		Id: id,
	}, nil
}
