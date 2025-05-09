package interceptors

import (
	"context"
	"github.com/peterouob/golang_template/utils"
	"google.golang.org/grpc/metadata"
	"time"

	"google.golang.org/grpc"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		md, _ := metadata.FromIncomingContext(ctx)
		resp, err := handler(ctx, req)
		utils.Logf("Method: %s, Time: %v,Metadata: %v", info.FullMethod, time.Since(start), md)
		return resp, err
	}
}
