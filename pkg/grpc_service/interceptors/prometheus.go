package interceptors

import (
	"context"
	promsever "github.com/peterouob/golang_template/pkg/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

func PromInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	s, m := parseInfo(info.FullMethod)
	start := time.Now()
	resp, err := handler(ctx, req)
	code := status.Code(err).String()

	defer func() {
		promsever.ReqMetrics.WithLabelValues(s, m, code).Inc()
		promsever.HisMetrics.WithLabelValues(s, m).Observe(time.Since(start).Seconds())
	}()
	return resp, err
}

func parseInfo(info string) (string, string) {
	s := strings.Split(info, "/")
	return s[1], s[2]
}
