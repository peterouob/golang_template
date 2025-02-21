package main

import (
	"flag"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/pkg/grpc_service/interceptors"
	"github.com/peterouob/golang_template/server"
	"github.com/peterouob/golang_template/tools"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

var (
	port = flag.Int("port", 8082, "grpc service port")
)

func init() {
	tools.InitLogger()
	configs.InitViper()
}

func main() {
	flag.Parse()
	go func() {
		http.Handle("/", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":9092", nil))
	}()

	servers := []server.GrpcServer{
		server.RegisterUserService("echo", nil, nil),
		server.RegisterUserService("login", nil, nil),
		server.RegisterUserService("jwt", []grpc.UnaryServerInterceptor{interceptors.TokenInterceptors}, nil),
		server.RegisterUserService("auth", nil, nil),
	}
	ports := []int{8081, 8082, 8083, 8084}
	for i, gserver := range servers {
		go gserver.InitServer(ports[i])
	}
	//server.GrpcGatewayServer(*port)
	select {}
}
