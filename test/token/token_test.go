package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/pkg/verify"
	"github.com/peterouob/golang_template/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	userId int64
	token  *verify.Token
)

func setup() {
	utils.InitLogger()
	configs.InitViper()
	userId = 123
	fmt.Printf("\033[1;33m%s\033[0m", "> Setup completed\n")
	token = verify.NewToken(userId)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestNewToken(t *testing.T) {
	token = verify.NewToken(int64(123456))
	assert.NotNil(t, token)
	assert.NotEqual(t, userId, token.UserId)
	assert.NotEmpty(t, token.Token.AccessUuid)
	assert.NotEmpty(t, token.Token.RefreshUuid)
	assert.Greater(t, token.Token.AtExpires, time.Now().Unix())
	assert.Greater(t, token.Token.RefreshAtExpires, time.Now().Unix())
	t.Logf("token: %v", token)
	t.Logf("token.Token: %v", token.Token)
}

func TestCreateToken(t *testing.T) {
	token := verify.NewToken(userId)
	assert.Equal(t, token.AccessToken, "")
	token.CreateToken()
	assert.NotEqual(t, token.AccessToken, "")

	parse, err := jwt.Parse(token.AccessToken, func(token *jwt.Token) (any, error) {
		var t []byte
		t = fmt.Appendf(t, fmt.Sprintf("%s", verify.TokenKey.Load().(string)))
		return t, nil
	})
	claims := parse.Claims.(jwt.MapClaims)
	t.Logf("parse: %v", claims["userId"])
	assert.NoError(t, err)
	assert.NotNil(t, parse)
}
