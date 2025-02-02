package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func GrpcGatewayServer(port int) {
	localIP := tools.GetLocalIP()
	ctxr := context.Background()
	ctx, cancel := context.WithCancel(ctxr)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := protobuf.RegisterEchoHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", localIP, port), opts)
	tools.HandelError("register grpc gateway server", err)
	tools.Log("register grpc gateway server success")
	tools.Log(fmt.Sprintf("grpc gateway server listening on port %s:%d", localIP, 7111))
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", localIP, 7111), mux)
}
