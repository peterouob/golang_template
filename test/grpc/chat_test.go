package grpc

import (
	"context"
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"testing"
)

func testBroadCast(t *testing.T, token string) {
	conn, err := grpc.NewClient("192.168.0.100:7082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	client := protobuf.NewChatClient(conn)
	assert.NotNil(t, client)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", fmt.Sprintf("Bearer %s", token))

	log.Println("THIS IS BROADCAST")

	_, err = client.BroadCast(ctx)
	assert.NoError(t, err)

}
