package server

import (
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/grpc_service/interceptors"
	grpcserver "github.com/peterouob/golang_template/pkg/grpc_service/server"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
)

type UserServiceSever struct {
	BaseServer
}

var userService = map[string]func([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor) *UserServiceSever{
	"echo":  newEchoService,
	"login": newLoginService,
	"jwt":   newjwtService,
}

func RegisterUserService(serviceName string, exUnIn []grpc.UnaryServerInterceptor, exStIn []grpc.StreamServerInterceptor) *UserServiceSever {
	if u, ok := userService[serviceName]; ok {
		return u(exUnIn, exStIn)
	}
	tools.ErrorMsg("No found service " + serviceName)
	return nil
}

func NewUserService(name string, regFunc func(server *grpc.Server), extUnIn []grpc.UnaryServerInterceptor, exStIn []grpc.StreamServerInterceptor) *UserServiceSever {
	baseInterceptor := []grpc.UnaryServerInterceptor{
		interceptors.LoggingInterceptor,
	}
	baseInterceptor = append(baseInterceptor, extUnIn...)
	server := &UserServiceSever{
		BaseServer{
			ServiceName:  name,
			RegisterFunc: regFunc,
		},
	}
	server.RegisterUnInterceptors(baseInterceptor...)
	server.RegisterStreamInterceptors(exStIn...)
	return server
}

func newEchoService(exUnIn []grpc.UnaryServerInterceptor, exStIn []grpc.StreamServerInterceptor) *UserServiceSever {
	return NewUserService("echo", func(server *grpc.Server) {
		echo := grpcserver.NewEchoServer()
		protobuf.RegisterEchoServer(server, echo)
		tools.Log("register echo service success")
	}, exUnIn, exStIn)
}

func newLoginService(exUnIn []grpc.UnaryServerInterceptor, exStIn []grpc.StreamServerInterceptor) *UserServiceSever {
	return NewUserService("echo", func(server *grpc.Server) {
		s := grpcserver.NewLoginServer()
		protobuf.RegisterUserServer(server, s)
		tools.Log("register echo service success")
	}, exUnIn, exStIn)
}

func newjwtService(exUnIn []grpc.UnaryServerInterceptor, exStIn []grpc.StreamServerInterceptor) *UserServiceSever {
	return NewUserService("echo", func(server *grpc.Server) {
		s := grpcserver.NewTokenTestServer()
		protobuf.RegisterUserServer(server, s)
		tools.Log("register echo service success")
	}, exUnIn, exStIn)
}
