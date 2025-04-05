package server

import (
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/grpc/interceptors"
	"github.com/peterouob/golang_template/pkg/grpc/server/im"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type IMServiceServer struct {
	BaseServer
}

var imService = map[string]func() *IMServiceServer{
	"broadcast": newBroadcast,
}

func RegisterIMService(serviceName string) *IMServiceServer {
	i, ok := imService[serviceName]
	if !ok {
		tools.ErrorMsgF("error in not found service name %s", serviceName)
	}
	return i()
}

func newIMService(name string, regFunc func(server *grpc.Server), exStIn grpc.StreamServerInterceptor) *IMServiceServer {
	baseInterceptor := []grpc.StreamServerInterceptor{
		interceptors.TokenStreamInterceptor,
	}
	if exStIn != nil {
		baseInterceptor = append(baseInterceptor, exStIn)
	}
	server := &IMServiceServer{
		BaseServer{
			ServiceName:  name,
			RegisterFunc: regFunc,
		},
	}
	server.RegisterStreamInterceptors(baseInterceptor...)
	return server
}

func newBroadcast() *IMServiceServer {
	return newIMService("broadcast", func(server *grpc.Server) {
		s := im.NewBroadCastServer()
		protobuf.RegisterChatServer(server, s)
		reflection.Register(server)
	}, nil)
}
