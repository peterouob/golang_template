package etcd

import (
	"context"
	"github.com/peterouob/golang_template/api/protobuf"
	grpcclient "github.com/peterouob/golang_template/pkg/grpc_service/client"
	"github.com/peterouob/golang_template/tools"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	tools.InitLogger()
}

func TestEchoEtcd(t *testing.T) {
	client := grpcclient.EchoClient()
	assert.NotEqual(t, nil, client)
	ctx := context.Background()
	resp, err := client.Echo(ctx, &protobuf.EchoRequest{
		Name: "hello",
	})
	assert.NoError(t, err)
	assert.Equal(t, "Server say hello to hello", resp.GetName())
}
