package repository

import (
	"context"
	"github.com/peterouob/golang_template/tools"
	"github.com/redis/go-redis/v9"
	"time"
)

type TokenRepo struct {
	rdb *redis.Client
}

var tokenRepo *TokenRepo

func GetTokenRepo() *TokenRepo {
	return tokenRepo
}

func NewTokenRepo(rdb *redis.Client) *TokenRepo {
	tools.Log("new token repo ...")
	t := &TokenRepo{rdb: rdb}
	tokenRepo = t
	return t
}

func (t *TokenRepo) SaveToken(ctx context.Context, id, key, value string, exp int64) {
	refreshTime := time.Until(time.Unix(exp, 0))
	ref := map[string]interface{}{
		"user_id": id,
		"token":   value,
		"exp":     exp,
		"create":  time.Now().Format("2006-01-02 15:04:05"),
	}

	err := t.rdb.HSet(ctx, key, ref).Err()
	tools.HandelError("error in save refresh token in redis HSet", err)
	err = t.rdb.Expire(ctx, key, refreshTime).Err()
	tools.HandelError("error in save refresh token in redis Expire", err)
}
