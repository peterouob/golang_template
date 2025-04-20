package etcd

//
//import (
//	"context"
//	"fmt"
//	"github.com/peterouob/golang_template/api/protobuf"
//	"github.com/peterouob/golang_template/configs"
//	grpcclient "github.com/peterouob/golang_template/pkg/grpc/client"
//	"github.com/peterouob/golang_template/utils"
//	"github.com/stretchr/testify/assert"
//	"google.golang.org/grpc/metadata"
//	"testing"
//)
//
//func init() {
//	utils.InitLogger()
//}
//
//func TestEchoEtcd(t *testing.T) {
//	cfg := &configs.EtcdGrpcCfg{}
//	cfg.SetPoolSize(4)
//	cfg.SetEndPoints([]string{"127.0.0.1:2379"})
//	cfg.SetServiceName("echo")
//	client, err := grpcclient.GetGRPCClient(cfg, "echo")
//	assert.NotNil(t, client)
//	assert.NoError(t, err)
//	resp, _ := client.(protobuf.EchoClient).Echo(context.Background(), &protobuf.EchoRequest{Name: "Hello"})
//	assert.NotNil(t, resp)
//}
//
//func TestLoginServer(t *testing.T) {
//	cfg := &configs.EtcdGrpcCfg{}
//	cfg.SetPoolSize(4)
//	cfg.SetEndPoints([]string{"127.0.0.1:2379"})
//	cfg.SetServiceName("login")
//	client, err := grpcclient.GetGRPCClient(cfg, "login")
//	assert.NoError(t, err)
//	r, err := client.(protobuf.UserClient).LoginUser(context.Background(), &protobuf.LoginUserRequest{
//		Email:    "admin",
//		Password: "admin",
//	})
//	assert.NotNil(t, r)
//	testToken(t, r.AccessToken)
//}
//
//func testToken(t *testing.T, token string) {
//	cfg := &configs.EtcdGrpcCfg{}
//	cfg.SetPoolSize(4)
//	cfg.SetEndPoints([]string{"127.0.0.1:2379"})
//	cfg.SetServiceName("token")
//	client, err := grpcclient.GetGRPCClient(cfg, "token")
//	assert.NoError(t, err)
//	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", fmt.Sprintf("Bearer %s", token))
//	r, _ := client.(protobuf.UserClient).TokenTest(ctx, &protobuf.TokenTestRequest{})
//	t.Logf("%v", r)
//}
