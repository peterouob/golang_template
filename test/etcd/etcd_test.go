package etcd

import (
	"context"
	"fmt"
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/configs"
	grpcclient "github.com/peterouob/golang_template/pkg/grpc_service/client"
	"github.com/peterouob/golang_template/tools"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
)

func init() {
	tools.InitLogger()
}

func TestEchoEtcd(t *testing.T) {
	cfg := &configs.EtcdGrpcCfg{}
	cfg.SetPoolSize(4)
	cfg.SetEndPoints([]string{"127.0.0.1:2379"})
	cfg.SetServiceName("echo_service")
	client, err := grpcclient.GetGRPCClient(cfg, "echo_service")
	assert.NotNil(t, client)
	assert.NoError(t, err)
	resp, err := client.(protobuf.EchoClient).Echo(context.Background(), &protobuf.EchoRequest{Name: "Hello"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestLoginServer(t *testing.T) {
	cfg := &configs.EtcdGrpcCfg{}
	cfg.SetPoolSize(4)
	cfg.SetEndPoints([]string{"127.0.0.1:2379"})
	cfg.SetServiceName("login_service")
	client, err := grpcclient.GetGRPCClient(cfg, "login_service")
	assert.NoError(t, err)
	r, err := client.(protobuf.UserClient).LoginUser(context.Background(), &protobuf.LoginUserRequest{
		Email:    "admin",
		Password: "admin",
	})
	assert.NotNil(t, r)
	testToken(t, r.AccessToken)
}

func testToken(t *testing.T, token string) {
	cfg := &configs.EtcdGrpcCfg{}
	cfg.SetPoolSize(4)
	cfg.SetEndPoints([]string{"127.0.0.1:2379"})
	cfg.SetServiceName("tokentest_service")
	client, err := grpcclient.GetGRPCClient(cfg, "tokentest_service")
	assert.NoError(t, err)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", fmt.Sprintf("Bearer %s", token))
	r, err := client.(protobuf.UserClient).TokenTest(ctx, &protobuf.TokenTestRequest{})
	t.Logf("%v", r)
}
