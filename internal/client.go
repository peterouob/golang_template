package main

import (
	"bufio"
	"context"
	"fmt"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	reg := prometheus.NewRegistry()
	grpcMetric := grpcprom.NewClientMetrics()
	reg.MustRegister(grpcMetric)

	conn, err := grpc.NewClient("192.168.0.101:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpcMetric.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcMetric.StreamClientInterceptor()))

	if err != nil {
		panic(err)
	}

	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", 9094)}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	c := protobuf.NewEchoClient(conn)
	go func() {
		for {
			_, err := c.Echo(context.Background(), &protobuf.EchoRequest{Name: "Test"})
			if err != nil {
				log.Printf("Calling the SayHello method unsuccessfully. ErrorInfo: %+v", err)
				log.Printf("You should to stop the process")
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("You can press n or N to stop the process of client")
	for scanner.Scan() {
		if strings.ToLower(scanner.Text()) == "n" {
			os.Exit(0)
		}
	}
}
