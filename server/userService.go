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

var userService = map[string]func() *UserServiceSever{
	"echo":  newEchoService,
	"login": newLoginService,
	"jwt":   newJwtService,
	"auth":  newAuthService,
}

func RegisterUserService(serviceName string) *UserServiceSever {
	u, ok := userService[serviceName]
	if !ok {
		tools.ErrorMsgF("error in not found service name %s", serviceName)
		return newNotFoundServer(serviceName)
	}
	return u()
}

func newNotFoundServer(name string) *UserServiceSever {
	return newUserService(name, func(s *grpc.Server) {
		protobuf.RegisterNotFoundServer(s, &protobuf.UnimplementedNotFoundServer{})
	}, nil, nil)
}

func newUserService(name string, regFunc func(server *grpc.Server), extUnIn grpc.UnaryServerInterceptor, exStIn grpc.StreamServerInterceptor) *UserServiceSever {
	baseInterceptor := []grpc.UnaryServerInterceptor{
		interceptors.LoggingInterceptor(),
	}
	if extUnIn != nil {
		baseInterceptor = append(baseInterceptor, extUnIn)
	}
	server := &UserServiceSever{
		BaseServer{
			ServiceName:  name,
			RegisterFunc: regFunc,
		},
	}
	server.RegisterUnInterceptors(baseInterceptor...)
	if exStIn != nil {
		server.RegisterStreamInterceptors(exStIn)
	}
	return server
}

func newEchoService() *UserServiceSever {
	return newUserService("echo", func(server *grpc.Server) {
		echo := grpcserver.NewEchoServer()
		protobuf.RegisterEchoServer(server, echo)
		tools.Log("register echo service success")
	}, nil, nil)
}

func newLoginService() *UserServiceSever {
	return newUserService("login", func(server *grpc.Server) {
		s := grpcserver.NewLoginServer()
		protobuf.RegisterUserServer(server, s)
		tools.Log("register login service success")
	}, nil, nil)
}

func newJwtService() *UserServiceSever {
	return newUserService("token", func(server *grpc.Server) {
		s := grpcserver.NewTokenTestServer()
		protobuf.RegisterUserServer(server, s)
		tools.Log("register jwt test service success")
	}, interceptors.TokenInterceptors(), nil)
}

func newAuthService() *UserServiceSever {
	return newUserService("auth", func(server *grpc.Server) {
		s := grpcserver.NewTokenValidServer()
		protobuf.RegisterUserServer(server, s)
		tools.Log("register token valid service success")
	}, nil, nil)
}
