package interceptors

import (
	"context"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

var cfg = &configs.EtcdGrpcCfg{}

func init() {
	cfg = &configs.EtcdGrpcCfg{}
	cfg.ServiceName = "auth"
	cfg.SetEndPoints([]string{"127.0.0.1:2379"})
	cfg.SetPoolSize(10)
}

func TokenInterceptors() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		utils.Log("start unary interceptor for token valid ...")
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
		}
		tokenString, err := extractToken(md)
		utils.HandelError("error in interceptor", err)

		conn, err := grpc.NewClient(":8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			utils.Error("error in grpc Client", err)
		}
		c := protobuf.NewUserClient(conn)
		res, err := c.TokenValid(ctx, &protobuf.TokenValidRequest{
			Token: tokenString,
		})
		if err != nil || !res.Valid {
			utils.HandelError("error in interceptor for valid token", err)
		}

		ctx = context.WithValue(ctx, "uid", res.Id)
		resp, err := handler(ctx, req)
		return resp, err
	}
}

func extractToken(md metadata.MD) (string, error) {
	authHeader, ok := md["authorization"]
	if len(authHeader) == 0 || !ok {
		return "", status.Errorf(codes.Unauthenticated, "missing authHeader from metadata")
	}
	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", status.Errorf(codes.Unauthenticated, "invalid auth header format")
	}
	return parts[1], nil
}
