package grpc

import (
	"context"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestEchoServer(t *testing.T) {
	conn, err := grpc.NewClient("192.168.0.100:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer func() {
		err = conn.Close()
		assert.NoError(t, err)
	}()

	name := "peter"
	c := protobuf.NewEchoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Echo(ctx, &protobuf.EchoRequest{Name: name})
	assert.NoError(t, err)
	t.Logf("r:%s", r.GetName())
}
