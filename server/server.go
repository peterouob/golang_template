package server

import (
	"fmt"
	etcdregister "github.com/peterouob/golang_template/pkg/etcd"
	in "github.com/peterouob/golang_template/pkg/grpc/interceptors"
	promsever "github.com/peterouob/golang_template/pkg/prometheus"
	"github.com/peterouob/golang_template/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

type GrpcServer interface {
	InitServer(port string) <-chan struct{}
}

type BaseServerInterface interface {
	RegisterUnInterceptors(interceptors ...grpc.UnaryServerInterceptor)
	RegisterStreamInterceptors(interceptors ...grpc.StreamServerInterceptor)
	registerInterceptors() (opts []grpc.ServerOption)
	InitServer(port string) <-chan struct{}
}

type BaseServer struct {
	ServiceName        string
	RegisterFunc       func(*grpc.Server)
	interceptors       []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	etcdClient         *etcdregister.EtcdRegister
}

var _ BaseServerInterface = (*BaseServer)(nil)

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

func (b *BaseServer) InitServer(port string) <-chan struct{} {
	utils.Log(fmt.Sprintf("Starting gRPC server [%s] on port %s ...", b.ServiceName, port))
	ready := make(chan struct{})
	addr := utils.FormatIP(port)
	go func() {
		lis, err := net.Listen("tcp", addr)
		utils.HandelError("error in listen addr", err)

		m := promsever.InitPrometheus()
		b.RegisterUnInterceptors(in.PromInterceptor(m))

		opts := b.registerInterceptors()
		opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     5 * time.Minute,
			MaxConnectionAge:      10 * time.Minute,
			MaxConnectionAgeGrace: 5 * time.Minute,
			Time:                  2 * time.Minute,
			Timeout:               20 * time.Second,
		}), grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             30 * time.Second,
			PermitWithoutStream: true,
		}))

		s := grpc.NewServer(opts...)

		if b.RegisterFunc == nil {
			utils.ErrorMsg("have not fund the register service")
		}

		b.RegisterFunc(s)

		//TODO:Rebuild Etcd service
		//etcd := etcdregister.NewEtcdRegister([]string{"127.0.0.1:2379"}, 3)
		//etcd.Register(b.ServiceName, addr)

		close(ready)
		err = s.Serve(lis)
		//b.listenExit(addr, s)
		utils.Error("error in start grpc server", err)
	}()
	return ready
}

//func (b *BaseServer) listenExit(addr string, s *grpc.Server) {
//	go func() {
//		c := make(chan os.Signal, 1)
//		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
//		sig := <-c
//		utils.Log(fmt.Sprintf("receive a signal %s", sig.String()))
//		s.GracefulStop()
//		b.cleanup(addr)
//	}()
//}
//
//func (b *BaseServer) cleanup(addr string) {
//	if b.etcdClient != nil {
//		b.etcdClient.UnRegister(b.ServiceName, addr)
//		utils.Log("Etcd unregistered successfully.")
//	}
//	utils.Log("Server shutdown finished.")
//}
