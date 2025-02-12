package interceptors

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
)

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	fmt.Printf("[Logging] gRPC method: %s\n", info.FullMethod)
	resp, err := handler(ctx, req)
	return resp, err
}
