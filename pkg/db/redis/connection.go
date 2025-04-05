package rdb

import (
	"context"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/tools"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:            configs.Config.GetString("redis.addr"),
		DB:              configs.Config.GetInt("redis.db"),
		Password:        configs.Config.GetString("redis.password"),
		MaxRetryBackoff: 5 * time.Minute,
		PoolSize:        10,
	})
	status := rdb.Ping(context.Background())
	tools.HandelError("error in ping redis", status.Err())
	return rdb
}
