package grpc_service

import (
	"fmt"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"time"
)

type PoolConn struct {
	conn     *grpc.ClientConn
	dealTime time.Time
}

func (pc *PoolConn) Refresh(cfg configs.ClientConfig, opts ...grpc.DialOption) {
	if pc == nil {
		return
	}
	if pc.conn != nil {
		err := pc.conn.Close()
		tools.HandelError("close pool connect conn", err)
	}
	conn, err := grpc.NewClient(cfg.ServerAddr, opts...)
	tools.HandelError(fmt.Sprintf("connect to %s failed", cfg.ServerAddr), err)
	pc.conn = conn
	pc.dealTime = time.Now().Add(cfg.GenLifeTime())
}

func (pc *PoolConn) ShouldRefresh() bool {
	if pc.conn == nil || !isConnectionHealthy(pc.conn) || time.Now().After(pc.dealTime) {
		return true
	}
	return false
}
