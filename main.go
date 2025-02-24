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
	port    = flag.Int("port", 8082, "grpc service port")
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
		server.RegisterUserService("echo"),
		server.RegisterUserService("login"),
		server.RegisterUserService("jwt"),
		server.RegisterUserService("auth"),
		server.RegisterUserService("register"),
	}
	ports := []int{8081, 8082, 8083, 8084, 8085}
	for i, gserver := range servers {
		go gserver.InitServer(ports[i])
	}
	//server.GrpcGatewayServer(*port)
	select {}
}
