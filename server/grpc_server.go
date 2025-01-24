package server

import (
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/etcd/server"
	"github.com/peterouob/golang_template/pkg/grpc_service"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"net"
	"time"
)

func InitGrpcServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	tools.HandelError("listen errors ", err)
	s := grpc.NewServer()
	echo := grpc_service.NewEchoServer()
	protobuf.RegisterEchoServer(s, echo)

	localIP := tools.GetLocalIP()

	var heart int64 = 3
	etcd := server.RegisterETCD([]string{"127.0.0.1:2379"}, heart)
	leaseId := etcd.Register("echo_service", fmt.Sprintf("%s:%d", localIP, port), 0)

	go func() {
		for {
			etcd.Register("echo_service", fmt.Sprintf("%s:%d", localIP, port), leaseId)
			time.Sleep(time.Duration(heart)*time.Second - 100*time.Millisecond)
		}
	}()

	err = s.Serve(lis)
	tools.HandelError("grpc server start error", err)
}
