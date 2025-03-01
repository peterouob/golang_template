package interceptors

import (
	"context"
	"errors"
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
	ctx context.Context
}

func (w *wrappedStream) RecvMsg(m any) error {
	tools.Logf("Receive a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m any) error {
	tools.Logf("Send a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func newWrappedStream(s grpc.ServerStream, ctx context.Context) grpc.ServerStream {
	return &wrappedStream{s, ctx}
}

func TokenStreamInterceptor(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	tools.Log("TokenStreamInterceptor is triggered")

	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.InvalidArgument, "missing metadata")
	}
	tokenString, err := extractToken(md)
	c, err := grpcclient.GetGRPCClient(cfg, "auth")
	tools.HandelError("error in interceptor for get grpc client", err)

	tctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := c.(protobuf.UserClient).TokenValid(tctx, &protobuf.TokenValidRequest{
		Token: tokenString,
	})
	ctx = context.WithValue(ctx, "id", res.Id)

	if err == nil && res.Valid == true {
		err = handler(srv, newWrappedStream(ss, ctx))
		tools.Log("Handler execution finished")
	} else if res.Valid != true {
		return status.Error(codes.PermissionDenied, "permission denied")
	}

	return errors.New("error in interceptor")
}
