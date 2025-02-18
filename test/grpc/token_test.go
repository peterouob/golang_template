package grpc

import (
	"context"
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func TestLoginServer(t *testing.T) {
	for {
		conn, err := grpc.NewClient("192.168.0.100:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
		assert.NoError(t, err)
		defer func() {
			err = conn.Close()
			assert.NoError(t, err)
		}()

		c := protobuf.NewUserClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
		defer cancel()

		r, err := c.LoginUser(ctx, &protobuf.LoginUserRequest{
			Email:    "admin",
			Password: "admin",
		})
		assert.NoError(t, err)
		t.Logf("Access Token :%s", r.AccessToken)
		t.Logf("Refresh Token :%s", r.RefreshToken)
		testToken(t, r.AccessToken)
		time.Sleep(1 * time.Second)
	}
}

func testToken(t *testing.T, token string) {
	conn, err := grpc.NewClient("192.168.0.100:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer func() {
		err = conn.Close()
		assert.NoError(t, err)
	}()
	c := protobuf.NewUserClient(conn)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", fmt.Sprintf("Bearer %s", token))
	_, err = c.TokenTest(ctx, &protobuf.TokenTestRequest{})
	assert.NoError(t, err)
}
