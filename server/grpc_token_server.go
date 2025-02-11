package server

import (
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/grpc_service/interceptors"
	grpcserver "github.com/peterouob/golang_template/pkg/grpc_service/server"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"net"
)

func InitTokenServer() {
	tools.Log("start grpc token server ...")
	localIp := tools.GetLocalIP()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIp, 8086))
	tools.HandelError("listen errors ", err)

	auth := interceptors.NewTokenInterceptor()

	s := grpc.NewServer(grpc.UnaryInterceptor(auth.UnaryServerInterceptor()))

	tokenServer := grpcserver.NewTokenTestServer()
	protobuf.RegisterUserServer(s, tokenServer)
	err = s.Serve(lis)
	tools.HandelError("login serve errors ", err)
}
