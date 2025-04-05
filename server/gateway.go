package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/configs"
	grpcclient "github.com/peterouob/golang_template/pkg/grpc/client"
	"github.com/peterouob/golang_template/tools"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type GatewayConfig struct {
	GatewayAddr string
	ServiceName string
	Port        int
	Cfg         *configs.EtcdGrpcCfg
	sync.RWMutex
}

const p int = 30001

func NewGatewayConfig(serviceName string, port int) *GatewayConfig {
	cfg := &configs.EtcdGrpcCfg{}
	cfg.SetPoolSize(10)
	cfg.SetEndPoints([]string{"127.0.0.1:2379"})

	return &GatewayConfig{
		Port:        port,
		ServiceName: serviceName,
		Cfg:         cfg,
	}
}

func (gw *GatewayConfig) StartGateway(wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(tools.Matcher))

	gw.Cfg.SetServiceName(gw.ServiceName)
	conn := grpcclient.GetGRPCUserClient(gw.Cfg)
	if err := pb.RegisterUserHandlerClient(ctx, mux, conn); err != nil {
		tools.ErrorMsg(fmt.Sprintf("Failed to register handler for %s: %v", gw.ServiceName, err))
	} else {
		//tools.Logf("Successfully registered handler for %s", gw.ServiceName)
	}

	handler := tools.Cors(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", tools.GetLocalIP(), gw.Port),
		Handler: handler,
	}

	go func() {
		tools.Logf("Start [%s] gRPC Gateway on port %d", gw.ServiceName, gw.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			tools.HandelError(fmt.Sprintf("Error starting gRPC Gateway [%s] on port %d", gw.ServiceName, gw.Port), err)
		}
	}()

	<-ctx.Done()
	tools.Logf("Shutting down gRPC Gateway: %s", gw.ServiceName)

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutDownCtx); err != nil {
		tools.HandelError(fmt.Sprintf("error in shutdown web server [%s]", gw.ServiceName), err)
	}
	tools.Logf("gRPC Gateway [%s] shutdown completed", gw.ServiceName)
}
