package grpc

import (
	"context"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestRegister(t *testing.T) {
	conn, err := grpc.NewClient("0.0.0.0:8085",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	defer func() {
		err = conn.Close()
		assert.Nil(t, err)
	}()

	c := protobuf.NewUserClient(conn)
	_, err = c.RegisterUser(context.Background(), &protobuf.RegisterUserRequest{
		Name:     "admin121233",
		Password: "12345613321322",
		Email:    "admin@a111231232dmin.com",
	})
	assert.Nil(t, err)
}
