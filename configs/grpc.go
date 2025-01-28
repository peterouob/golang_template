package configs

import (
	"math"
	"math/rand/v2"
	"time"
)

type ClientConfig struct {
	ServerAddr        string
	PoolSize          int
	LifeTime          time.Duration
	LifeTimeDeviation time.Duration
}

//採用builder策略

func (cc *ClientConfig) SetServerAddr(addr string) *ClientConfig {
	cc.ServerAddr = addr
	return cc
}

func (cc *ClientConfig) SetPoolSize(size int) *ClientConfig {
	cc.PoolSize = size
	return cc
}

func (cc *ClientConfig) SetLifeTime(lifeTime time.Duration) *ClientConfig {
	cc.LifeTime = lifeTime
	return cc
}

func (cc *ClientConfig) SetLifeTimeDeviation(lifeTime time.Duration) *ClientConfig {
	cc.LifeTimeDeviation = lifeTime
	return cc
}

func (cc *ClientConfig) GenLifeTime() time.Duration {
	if cc.LifeTime > 0 {
		random := time.Duration(int64(rand.Float64() * float64(cc.LifeTimeDeviation)))
		return cc.LifeTimeDeviation + random
	}
	return time.Duration(math.MaxInt64 - time.Now().Unix())
}
