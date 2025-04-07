package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/peterouob/golang_template/api/router"
	"github.com/peterouob/golang_template/configs"
	mdb "github.com/peterouob/golang_template/pkg/db/mysql"
	rdb "github.com/peterouob/golang_template/pkg/db/redis"
	"github.com/peterouob/golang_template/pkg/repository"
	"github.com/peterouob/golang_template/server"
	"github.com/peterouob/golang_template/utils"
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
	utils.InitLogger()
	configs.InitViper()
	mysqldb = mdb.InitMysql()
	redisdb = rdb.InitRedis()
}

func main() {
	flag.Parse()
	go func() {
		repository.NewUserRepo(mysqldb)
		repository.NewTokenRepo(redisdb)
		http.Handle("/", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":9092", nil))
	}()

	servers := []server.GrpcServer{
		server.RegisterUserService("login"),
		server.RegisterUserService("token"),
		server.RegisterUserService("auth"),
		server.RegisterUserService("register"),
		server.RegisterIMService("broadcast"),
	}
	ports := []string{"8082", "8083", "8084", "8085", "7082"}

	readies := make([]<-chan struct{}, len(servers))
	for i, s := range servers {
		readies[i] = s.InitServer(ports[i])
	}
	for _, ch := range readies {
		<-ch
	}

	r := gin.Default()
	router.InitRouter(r)
	if err := r.Run(":9093"); err != nil {
		utils.Error("error in open gin server", err)
	}
	//select {}
}
