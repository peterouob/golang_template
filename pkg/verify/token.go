package verify

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/peterouob/golang_template/configs"
	"github.com/peterouob/golang_template/tools"
	"sync/atomic"
	"time"
)

var (
	err        error
	TokenKey   atomic.Value
	RefreshKey atomic.Value
)

type Token struct {
	UserId       int64         `json:"user_id"`
	AccessId     string        `json:"access_id"`
	AccessToken  string        `json:"access_token"`
	RefreshId    string        `json:"refresh_id"`
	RefreshToken string        `json:"refresh_token"`
	Token        configs.Token `json:"token"`
}

func NewToken(id int64) *Token {
	TokenKey.Store(configs.Config.GetString("token.token_key"))
	RefreshKey.Store(configs.Config.GetString("token.refresh_key"))
	token := &configs.Token{}
	token.AccessUuid = uuid.NewString()
	token.RefreshUuid = uuid.NewString()
	token.AtExpires = time.Now().Add(time.Hour * 2).Unix()
	token.RefreshAtExpires = time.Now().Add(time.Hour * 24 * 7 * 2).Unix()
	return &Token{
		UserId: id,
		Token:  *token,
	}
}

func (t *Token) CreateToken() {
	claims := jwt.MapClaims{
		"access_id": t.Token.AccessUuid,
		"exp":       t.Token.AtExpires,
		"type":      "access",
		"userId":    t.UserId,
		"jti":       t.UserId,
		"iat":       time.Now().Unix(),
	}

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t.AccessToken, err = tk.SignedString([]byte(TokenKey.Load().(string)))
	tools.HandelError("create token error", err)
	t.AccessId = claims["access_id"].(string)
}
