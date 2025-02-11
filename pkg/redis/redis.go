package redis

import (
	"context"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/tools"
	"github.com/redis/go-redis/v9"
	"time"
)

var Rdb *redis.Client

func InitRedis() {
	//TODO:redis cluster connect now only one connect
	rdb := redis.NewClient(&redis.Options{
		Addr:            configs.Config.GetString("redis.addr"),
		DB:              configs.Config.GetInt("redis.db"),
		Password:        configs.Config.GetString("redis.password"),
		MaxRetryBackoff: 5 * time.Minute,
		PoolSize:        10,
	})
	status := rdb.Ping(context.Background())
	tools.HandelError("error in ping redis", status.Err())
	Rdb = rdb
}
