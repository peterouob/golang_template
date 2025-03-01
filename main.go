package main

import (
	"flag"
	"github.com/peterouob/golang_template/configs"
	mdb "github.com/peterouob/golang_template/pkg/db/mysql"
	rdb "github.com/peterouob/golang_template/pkg/db/redis"
	"github.com/peterouob/golang_template/pkg/repository"
	"github.com/peterouob/golang_template/server"
	"github.com/peterouob/golang_template/tools"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var (
	mysqldb *gorm.DB
	redisdb *redis.Client
)

func init() {
	tools.InitLogger()
	configs.InitViper()
	mysqldb = mdb.InitMysql()
	redisdb = rdb.InitRedis()
}

func main() {
	flag.Parse()
	go func() {
		http.Handle("/", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":9092", nil))
	}()
	go func() {
		repository.NewUserRepo(mysqldb)
		repository.NewTokenRepo(redisdb)
	}()

	servers := []server.GrpcServer{
		server.RegisterUserService("login"),
		server.RegisterUserService("token"),
		server.RegisterUserService("auth"),
		server.RegisterUserService("register"),
		server.RegisterIMService("broadcast"),
	}
	ports := []int{8082, 8083, 8084, 8085, 7082}

	readies := make([]<-chan struct{}, len(servers))
	for i, s := range servers {
		readies[i] = s.InitServer(ports[i])
	}
	for _, ch := range readies {
		<-ch
	}
	tools.Log("All gRPC services started. Starting gRPC Gateway...")

	//var wg sync.WaitGroup
	//gateways := []*server.GatewayConfig{
	//	server.NewGatewayConfig("login", 30001),
	//	server.NewGatewayConfig("register", 30002),
	//	server.NewGatewayConfig("token", 30003),
	//}
	//for _, gw := range gateways {
	//	wg.Add(1)
	//	go gw.StartGateway(&wg)
	//}
	select {}
}
