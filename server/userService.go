package server

import (
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/grpc/interceptors"
	"github.com/peterouob/golang_template/pkg/grpc/server/user"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UserServiceSever struct {
	BaseServer
}

var userService = map[string]func() *UserServiceSever{
	"login":    newLoginService,
	"token":    newJwtService,
	"auth":     newAuthService,
	"register": newRegisterService,
}

func RegisterUserService(serviceName string) *UserServiceSever {
	u, ok := userService[serviceName]
	if !ok {
		tools.ErrorMsgF("error in not found service name %s", serviceName)
	}
	return u()
}

func newUserService(name string, regFunc func(server *grpc.Server), extUnIn ...grpc.UnaryServerInterceptor) *UserServiceSever {
	baseInterceptor := []grpc.UnaryServerInterceptor{
		interceptors.LoggingInterceptor(),
	}
	if extUnIn != nil {
		baseInterceptor = append(baseInterceptor, extUnIn...)
	}
	server := &UserServiceSever{
		BaseServer{
			ServiceName:  name,
			RegisterFunc: regFunc,
		},
	}
	server.RegisterUnInterceptors(baseInterceptor...)

	return server
}

func newLoginService() *UserServiceSever {
	return newUserService("login", func(server *grpc.Server) {
		s := user.NewLoginServer()
		protobuf.RegisterUserServer(server, s)
		reflection.Register(server)
		tools.Log("register login service success")
	})
}

func newJwtService() *UserServiceSever {
	return newUserService("token", func(server *grpc.Server) {
		s := user.NewTokenTestServer()
		protobuf.RegisterUserServer(server, s)
		reflection.Register(server)
		tools.Log("register jwt test service success")
	}, interceptors.TokenInterceptors())
}

func newAuthService() *UserServiceSever {
	return newUserService("auth", func(server *grpc.Server) {
		s := user.NewTokenValidServer()
		protobuf.RegisterUserServer(server, s)
		reflection.Register(server)
		tools.Log("register token valid service success")
	})
}

func newRegisterService() *UserServiceSever {
	return newUserService("register", func(server *grpc.Server) {
		s := user.NewRegisterServer()
		protobuf.RegisterUserServer(server, s)
		reflection.Register(server)
		tools.Log("register user register service success")
	})
}
