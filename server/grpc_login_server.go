package server

import (
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	grpcserver "github.com/peterouob/golang_template/pkg/grpc_service/server"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"net"
)

func InitLoginServer() {
	tools.Log("start grpc login server ...")
	localIp := tools.GetLocalIP()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIp, 8085))
	tools.HandelError("listen errors ", err)
	s := grpc.NewServer()
	login := grpcserver.NewLoginServer()
	protobuf.RegisterUserServer(s, login)
	err = s.Serve(lis)
	tools.HandelError("login serve errors ", err)
}
