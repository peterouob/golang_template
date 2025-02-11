package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/pkg/verify"
	"github.com/peterouob/golang_template/tools"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewToken(t *testing.T) {
	tools.InitLogger()
	configs.InitViper()
	userId := int64(123)
	token := verify.NewToken(userId)
	assert.NotNil(t, token)
	assert.Equal(t, userId, token.UserId)
	assert.NotEmpty(t, token.Token.AccessUuid)
	assert.NotEmpty(t, token.Token.RefreshUuid)
	assert.Greater(t, token.Token.AtExpires, time.Now().Unix())
	assert.Greater(t, token.Token.RefreshAtExpires, time.Now().Unix())
	t.Logf("token: %v", token)
	t.Logf("token.Token: %v", token.Token)
}

func TestCreateToken(t *testing.T) {
	tools.InitLogger()
	configs.InitViper()
	userId := int64(123)
	token := verify.NewToken(userId)
	assert.Equal(t, token.AccessToken, "")
	token.CreateToken()
	assert.NotEqual(t, token.AccessToken, "")
	t.Log(verify.TokenKey.Load().(string))

	parse, err := jwt.Parse(token.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(verify.TokenKey.Load().(string)), nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, parse)
}
