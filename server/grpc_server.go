package server

import (
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/etcd/server"
	"github.com/peterouob/golang_template/pkg/grpc_service/server"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func InitGrpcServer(port int) {
	tools.Log("start grpc server ...")
	localIP := tools.GetLocalIP()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIP, port))
	tools.HandelError("listen errors ", err)
	s := grpc.NewServer()
	echo := grpcserver.NewEchoServer()
	protobuf.RegisterEchoServer(s, echo)

	var heart int64 = 3
	etcd := etcdservice.RegisterETCD([]string{"127.0.0.1:2379"}, heart)
	leaseId := etcd.Register("echo_service", fmt.Sprintf("%s:%d", localIP, port), 0)

	go func() {
		for {
			etcd.Register("echo_service", fmt.Sprintf("%s:%d", localIP, port), leaseId)
			time.Sleep(time.Duration(heart)*time.Second - 100*time.Millisecond)
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		tools.Log(fmt.Sprintf("receive a signal %s", sig.String()))
		etcd.UnRegister("echo_service", fmt.Sprintf("%s:%d", localIP, port))
		os.Exit(0)
	}()

	err = s.Serve(lis)
	tools.HandelError("start grpc server error", err,
		func(args ...interface{}) {
			etcd.UnRegister("echo_service", fmt.Sprintf("%s:%d", localIP, port))
		})
}
