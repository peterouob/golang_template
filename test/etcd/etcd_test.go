package etcd

import (
	"context"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/configs"
	grpcclient "github.com/peterouob/golang_template/pkg/grpc_service/client"
	"github.com/peterouob/golang_template/tools"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	tools.InitLogger()
}

func TestEchoEtcd(t *testing.T) {
	cfg := &configs.EtcdGrpcCfg{}
	//cfg.SetPoolSize(10)
	//cfg.SetEndPoints([]string{"127.0.0.1:2379"})
	//cfg.SetServiceName("echo_service")
	client, pool, err := grpcclient.EchoClient(cfg)
	assert.NotNil(t, client)
	assert.NotNil(t, pool)
	assert.NoError(t, err)
	resp, err := client.Echo(context.Background(), &protobuf.EchoRequest{Name: "Hello"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
