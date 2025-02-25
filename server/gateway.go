package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

var gatewayFactor = map[int]func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error{
	8081: pb.RegisterEchoHandlerFromEndpoint,
	8082: pb.RegisterUserHandlerFromEndpoint,
	8083: pb.RegisterUserHandlerFromEndpoint,
	8084: pb.RegisterUserHandlerFromEndpoint,
	8085: pb.RegisterUserHandlerFromEndpoint,
}

const p = 30001

func StartGateway(ports []int) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(tools.Matcher))
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials())}
	for _, port := range ports {
		f, ok := gatewayFactor[port]
		if !ok {
			tools.Logf("not found the service from port [%d] ", port)
			continue
		}
		addr := fmt.Sprintf("%s:%d", tools.GetLocalIP(), port)
		if err := f(ctx, mux, addr, opts); err != nil {
			tools.HandelError(fmt.Sprintf("Error registering gRPC Gateway on port %d", port), err)
		} else {
			tools.Logf("Gateway registered for port [%d] at %s", port, addr)
		}
	}
	handler := tools.Cors(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", tools.GetLocalIP(), p),
		Handler: handler,
	}

	go func() {
		tools.Logf("Start gRPC Gateway on port %d", p)
		if err := server.ListenAndServe(); err != nil {
			tools.HandelError(fmt.Sprintf("Error starting gRPC Gateway on port %d", p), err)
		}
	}()

	<-ctx.Done()
	tools.Log("Shutting down gRPC ...")

	shutDown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutDown); err != nil {
		tools.HandelError("error in shutdown web server by timeout", err)
	}
	tools.Log("gRPC Gateway shutdown completed")
}
