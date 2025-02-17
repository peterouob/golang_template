package main

//import (
//	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
//	"github.com/peterouob/golang_template/api/protobuf"
//	"github.com/peterouob/golang_template/pkg/grpc_service/interceptors"
//	grpcserver "github.com/peterouob/golang_template/pkg/grpc_service/server"
//	"github.com/peterouob/golang_template/tools"
//	"github.com/prometheus/client_golang/prometheus"
//	"google.golang.org/grpc"
//)
//
//type EchoServer struct {
//	BaseServer
//}
//
//var (
//	// Create a metrics registry.
//	reg = prometheus.NewRegistry()
//
//	// Create some standard server metrics.
//	grpcMetrics = grpcprom.NewServerMetrics()
//	// Create a customized counter metric.
//	customizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
//		Name: "demo_server_say_hello_method_handle_count",
//		Help: "Total number of RPCs handled on the server.",
//	}, []string{"name"})
//)
//
//func init() {
//	// Register standard server metrics and customized metrics to registry.
//	reg.MustRegister(grpcMetrics, customizedCounterMetric)
//	customizedCounterMetric.WithLabelValues("Test")
//}
//
//func RegisterEchoServer() *EchoServer {
//	echogrpc := &EchoServer{
//		BaseServer{
//			ServiceName: "echo_service",
//			RegisterFunc: func(server *grpc.Server) {
//				echo := grpcserver.NewEchoServer()
//				protobuf.RegisterEchoServer(server, echo)
//				grpcMetrics.InitializeMetrics(server)
//				tools.Log("register echo server success")
//			},
//		},
//	}
//	echogrpc.RegisterInterceptors(interceptors.LoggingInterceptor, grpcMetrics.UnaryServerInterceptor())
//	return echogrpc
//}
