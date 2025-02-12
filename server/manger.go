package server

import (
	"fmt"
	etcdregister "github.com/peterouob/golang_template/pkg/etcd"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"net"
)

type GrpcServer interface {
	InitServer(port int)
}

type BaseServer struct {
	ServiceName  string
	Registerfunc func(*grpc.Server)
	interceptors []grpc.UnaryServerInterceptor
}

func (b *BaseServer) RegisterInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	b.interceptors = append(b.interceptors, interceptors...)
}

func (b *BaseServer) InitServer(port int) {
	tools.Log(fmt.Sprintf("Starting gRPC server [%s] on port %d ...", b.ServiceName, port))
	addr := tools.FormatAddr(port)
	lis, err := net.Listen("tcp", addr)
	tools.HandelError("error in listen addr", err)

	var opts []grpc.ServerOption
	if len(b.interceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(b.interceptors...))
	}
	s := grpc.NewServer(opts...)
	if b.Registerfunc != nil {
		b.Registerfunc(s)
	} else {
		tools.Log("WARNING: RegisterFn is nil, no service registered")
	}

	etcd := etcdregister.NewEtcdRegister([]string{"127.0.0.1:2379"}, 3)
	etcd.Register(b.ServiceName, addr)
	etcd.ListenExit(b.ServiceName, addr)

	err = s.Serve(lis)
	tools.HandelError("start grpc server error", err,
		func(args ...interface{}) {
			etcd.UnRegister(b.ServiceName, addr)
		})
}
