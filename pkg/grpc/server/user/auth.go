package user

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/pkg/verify"
)

type TokenValid struct {
	protobuf.UnimplementedUserServer
}

func NewTokenValidServer() *TokenValid {
	return &TokenValid{}
}

func (auth TokenValid) TokenValid(ctx context.Context, req *protobuf.TokenValidRequest) (*protobuf.TokenValidResponse, error) {
	tokenString := req.GetToken()
	token := verify.TokenVerify(tokenString)

	if !token.Valid {
		return &protobuf.TokenValidResponse{
			Valid: false,
			Msg:   "Token is invalid",
		}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &protobuf.TokenValidResponse{
			Valid: false,
			Msg:   "Token claims are not valid",
		}, nil
	}
	userID := int64(claims["userId"].(float64))
	return &protobuf.TokenValidResponse{
		Valid: true,
		Id:    uint64(userID),
		Msg:   "Token is valid",
	}, nil
}
