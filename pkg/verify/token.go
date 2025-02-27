package verify

import (
	"context"
	"errors"
	"fmt"
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

// CreateToken  不存Redis單純驗證
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

// CreateRefreshToken TODO:存Redis並實現black list使用
func (t *Token) CreateRefreshToken() {
	claims := jwt.MapClaims{
		"refresh_id": t.Token.RefreshUuid,
		"exp":        t.Token.RefreshAtExpires,
		"type":       "refresh",
		"userId":     t.UserId,
		"jti":        t.UserId,
		"iat":        time.Now().Unix(),
	}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t.RefreshToken, err = tk.SignedString([]byte(fmt.Sprintf("%s%d", RefreshKey.Load().(string), t.UserId)))
	tools.HandelError("create refresh token error", err)
	t.RefreshId = claims["refresh_id"].(string)
}

func VerifyToken(tokenString string) *jwt.Token {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			tools.HandelError("parse token error type", err)
		}
		return []byte(TokenKey.Load().(string)), nil
	})
	tools.HandelError("parse token error", err)
	// TODO:count the fail and report to prometheus count
	switch {
	case token.Valid:
		tools.Log("valid success token")
	case errors.Is(err, jwt.ErrTokenMalformed):
		tools.Log("error in Malformed token type")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		tools.Log("error in token expired")
	default:
		tools.HandelError("couldn't handle this token", err)
	}
	return token
}

// SaveToken 等black list時使用
func SaveToken(ctx context.Context, id int64) (string, string) {
	token := NewToken(id)
	token.CreateToken()
	token.CreateRefreshToken()
	//exp := token.Token.GetRefreshAtExpires()
	//tokenRepo := repository.GetTokenRepo()
	//tokenRepo.SaveRefreshToken(ctx, strconv.FormatInt(id, 10), token.RefreshToken, exp)
	return token.AccessToken, token.RefreshToken
}
