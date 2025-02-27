package interceptors

import (
	"context"
	"github.com/peterouob/golang_template/api/protobuf"
	grpcclient "github.com/peterouob/golang_template/pkg/grpc/client"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m any) error {
	tools.Logf("Receive a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m any) error {
	tools.Logf("Send a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func TokenStreamInterceptor(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	tools.Log("TokenStreamInterceptor is triggered")

	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.InvalidArgument, "missing metadata")
	}
	tokenString, err := extractToken(md)
	tools.HandelError("error in interceptor", err)

	c, err := grpcclient.GetGRPCClient(cfg, "auth")
	tools.HandelError("error in interceptor for get grpc client", err)

	res, err := c.(protobuf.UserClient).TokenValid(ctx, &protobuf.TokenValidRequest{
		Token: tokenString,
	})
	if err != nil || !res.Valid {
		tools.HandelError("error in interceptor for valid token", err)
	}

	ctx = context.Background()

	err = handler(srv, newWrappedStream(ss))
	tools.Log("Handler execution finished")

	return err
}
