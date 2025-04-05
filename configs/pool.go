package configs

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"time"
)

const (
	DialTimeout = 5 * time.Second

	BackoffMaxDelay = 3 * time.Second

	KeepAliveTime = time.Duration(10) * time.Second

	KeepAliveTimeout = time.Duration(3) * time.Second

	InitialWindowSize = 1 << 30

	InitialConnWindowSize = 1 << 30

	MaxSendMsgSize = 4 << 30

	MaxRecvMsgSize = 4 << 30
)

type Option struct {
	Dial                 func(addr string) (*grpc.ClientConn, error)
	MaxIdle              int
	MaxActive            int
	MaxConcurrentStreams int
	Reuse                bool
}

var DefaultOption = Option{
	Dial:                 Dial,
	MaxIdle:              10,
	MaxActive:            10,
	MaxConcurrentStreams: 64,
	Reuse:                true,
}

func Dial(addr string) (*grpc.ClientConn, error) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	g, err := grpc.NewClient(addr, grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, errors.New(fmt.Sprint("failed to create gRPC client: ", err))
	}
	return g, nil
}
