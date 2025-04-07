package client

import (
	"context"
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/configs"
	mdb "github.com/peterouob/golang_template/pkg/db/mysql"
	"github.com/peterouob/golang_template/pkg/grpc/pool"
	"github.com/peterouob/golang_template/utils"
)

func LoginUserGrpc(addr string, c context.Context, model mdb.UserModel) *protobuf.LoginUserResponse {
	p := pool.New(fmt.Sprintf("192.168.0.100:%s", addr), configs.DefaultOption)
	conn, err := p.Get()
	if err != nil {
		utils.Error("get conn from pool error", err)
	}
	client := protobuf.NewUserClient(conn.Value())
	user := &protobuf.LoginUserRequest{
		Email:    model.Email,
		Password: model.Password,
	}
	resp, err := client.LoginUser(c, user)
	if err != nil {
		utils.Error("error in login user service", err)
	}
	return resp
}

func RegisterUser(addr string, c context.Context, model mdb.UserModel) *protobuf.RegisterUserResponse {
	p := pool.New(fmt.Sprintf("192.168.0.100:%s", addr), configs.DefaultOption)
	conn, err := p.Get()
	if err != nil {
		utils.Error("get conn from pool error", err)
	}
	client := protobuf.NewUserClient(conn.Value())
	user := &protobuf.RegisterUserRequest{
		Email:    model.Email,
		Name:     model.Name,
		Password: model.Password,
	}
	resp, err := client.RegisterUser(c, user)
	if err != nil {
		utils.Error("error in register user", err)
	}
	return resp
}
