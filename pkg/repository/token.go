package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/peterouob/golang_template/utils"
	"github.com/redis/go-redis/v9"
	"time"
)

type TokenRepo struct {
	rdb *redis.Client
}

var tokenRepo *TokenRepo

func NewTokenRepo(rdb *redis.Client) *TokenRepo {
	t := &TokenRepo{rdb: rdb}
	tokenRepo = t
	return t
}

func (t *TokenRepo) SaveRefreshToken(ctx context.Context, userId, value string, exp int64) {
	rttl := time.Until(time.Unix(exp, 0))
	userData := map[string]any{
		"token":  value,
		"exp":    exp,
		"create": time.Now().Format("2006-01-02 15:04:05"),
	}
	uBytes, err := json.Marshal(userData)
	utils.Error("json marshal fail", err)

	redisKey := fmt.Sprintf("refresh_%s", userId)

	err = t.rdb.HSet(ctx, redisKey, "black List", uBytes).Err()
	utils.Error("redis store hset fail", err)
	err = t.rdb.Expire(ctx, redisKey, rttl).Err()
	utils.HandelError("error in set expire for refresh token", err)
}

func (t *TokenRepo) GetRefreshTokenData(
	ctx context.Context,
	userId string,
) map[string]any {
	redisKey := fmt.Sprintf("user_refresh:%s", userId)
	dataBytes, err := t.rdb.HGet(ctx, redisKey, userId).Bytes()
	utils.HandelError("redis store hget data fail", err)
	var dataMap map[string]any
	err = json.Unmarshal(dataBytes, &dataMap)
	utils.HandelError("json unmarshal fail", err)
	return dataMap
}
