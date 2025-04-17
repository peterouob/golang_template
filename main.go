package main

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/golang_template/api/router"
	"github.com/peterouob/golang_template/configs"
	mdb "github.com/peterouob/golang_template/pkg/db/mysql"
	rdb "github.com/peterouob/golang_template/pkg/db/redis"
	"github.com/peterouob/golang_template/pkg/repository"
	"github.com/peterouob/golang_template/server"
	"github.com/peterouob/golang_template/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"sync"
)

func init() {
	utils.InitLogger()
	configs.InitViper()
	mysqldb := mdb.InitMysql()
	redisdb := rdb.InitRedis()
	repository.NewUserRepo(mysqldb)
	repository.NewTokenRepo(redisdb)
}

func main() {
	go startPrometheus()
	startRpcServer()
	startRouter()
}

func startRpcServer() {
	servers := []server.GrpcServer{
		server.RegisterUserService("login"),
		server.RegisterUserService("token"),
		server.RegisterUserService("auth"),
		server.RegisterUserService("register"),
		server.RegisterIMService("broadcast"),
	}
	ports := []string{"8082", "8083", "8084", "8085", "7082"}

	wg := sync.WaitGroup{}
	wg.Add(len(servers))
	for i, s := range servers {
		s.InitServer(ports[i])
		wg.Done()
	}
	wg.Wait()
}

func startRouter() {
	r := gin.Default()
	router.InitRouter(r)
	if err := r.Run(":9093"); err != nil {
		utils.Error("error in open gin server", err)
	}
}

func startPrometheus() {
	http.Handle("/", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9092", nil))
}
