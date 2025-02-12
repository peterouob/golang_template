package interceptors

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/peterouob/golang_template/pkg/verify"
	"github.com/peterouob/golang_template/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func TokenInterceptors(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	tools.Log("start unary interceptor ...")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	tokenString, err := extractToken(md)
	tools.HandelError("error in interceptor", err)
	token := verify.VerifyToken(tokenString)

	uid := int64(token.Claims.(jwt.MapClaims)["userId"].(float64))
	ctx = context.WithValue(ctx, "uid", uid)
	return handler(ctx, req)
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
