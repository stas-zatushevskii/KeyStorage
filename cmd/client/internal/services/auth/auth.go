package auth

import (
	"bytes"
	"client/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zalando/go-keyring"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func Get() (*Claims, error) {
	var claims Claims
	var t JWT

	// do request, get link from config
	buf := bytes.Buffer{} // fixme
	tokenString := strings.TrimSpace(string(buf.Bytes()))
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &claims)
	if err != nil {
		return nil, err
	}
	t.token = token
	t.expiration = claims.ExpiresAt.Time

	return &t, nil
}

func (c *JWT) Save() error {
	keyring.Set(config.App.JWTAccessToken, c)
	keyring
}

func (c *JWT) Expired() bool {
	return c.expiration.Sub(time.Now()) <= 0
}
