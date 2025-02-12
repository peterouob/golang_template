package grpcpool

import (
	"fmt"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"time"
)

type poolConn struct {
	conn     *grpc.ClientConn
	dealTime time.Time
}

func (pc *poolConn) Refresh(cfg configs.ClientConfig, opts ...grpc.DialOption) {
	if pc == nil {
		return
	}
	if pc.conn != nil {
		pc.conn.Close()
	}
	conn, err := grpc.NewClient(cfg.ServerAddr, opts...)
	tools.HandelError(fmt.Sprintf("connect to %s failed", cfg.ServerAddr), err)
	pc.conn = conn
	pc.dealTime = time.Now().Add(cfg.GetLifeTime())
}

func (pc *poolConn) ShouldRefresh() bool {
	if pc.conn == nil {
		tools.Log("refresh connection because conn is nil!")
		return true
	}

	if !isConnectionHealthy(pc.conn) {
		tools.Log("refresh connection because healthy!")
		return true
	}

	if time.Now().After(pc.dealTime) {
		tools.Log("refresh connection because dealTime!")
		return true
	}

	return false
}
