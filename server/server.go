package server

import (
	"fmt"
	etcdregister "github.com/peterouob/golang_template/pkg/etcd"
	in "github.com/peterouob/golang_template/pkg/grpc_service/interceptors"
	promsever "github.com/peterouob/golang_template/pkg/prometheus"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"net"
)

type GrpcServer interface {
	InitServer(port int)
}

type BaseServer struct {
	ServiceName        string
	RegisterFunc       func(*grpc.Server)
	interceptors       []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

func (b *BaseServer) RegisterUnInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	b.interceptors = append(b.interceptors, interceptors...)
}

func (b *BaseServer) RegisterStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	b.streamInterceptors = append(b.streamInterceptors, interceptors...)
}

func (b *BaseServer) registerInterceptors() (opts []grpc.ServerOption) {
	if len(b.interceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(b.interceptors...))
	}
	if len(b.streamInterceptors) > 0 {
		opts = append(opts, grpc.ChainStreamInterceptor(b.streamInterceptors...))
	}
	return
}

func (b *BaseServer) InitServer(port int) {
	tools.Log(fmt.Sprintf("Starting gRPC server [%s] on port %d ...", b.ServiceName, port))

	addr := tools.FormatAddr(port)
	lis, err := net.Listen("tcp", addr)
	tools.HandelError("error in listen addr", err)

	m := promsever.InitPrometheus()
	b.RegisterUnInterceptors(in.PromInterceptor(m))

	opts := b.registerInterceptors()
	s := grpc.NewServer(opts...)

	if b.RegisterFunc == nil {
		tools.ErrorMsg("have not fund the register service")
	}

	b.RegisterFunc(s)

	etcd := etcdregister.NewEtcdRegister([]string{"127.0.0.1:2379"}, 3)
	etcd.Register(b.ServiceName, addr)

	err = s.Serve(lis)
	tools.HandelError("start grpc server error", err,
		func(args ...interface{}) {
			etcd.UnRegister(b.ServiceName, addr)
		})
}
