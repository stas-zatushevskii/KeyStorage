package token

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	tokenJWT     string
	tokenRefresh string
	expiredAt    int64
}

func NewToken() *Token {
	return &Token{}
}

// SetJWTToken parse token from string and returns Token object with field expiresAt
func (t *Token) SetJWTToken(token string) error {
	var claims jwt.RegisteredClaims

	tokenString := strings.TrimSpace(token)
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, &claims)
	if err != nil {
		return err
	}
	t.tokenJWT = tokenString
	t.expiredAt = claims.ExpiresAt.Time.Unix()

	return nil
}

func (t *Token) GetJWTToken() string {
	return t.tokenJWT
}

func (t *Token) SetRefreshToken(token string) {
	t.tokenRefresh = token
}

func (t *Token) GetRefreshToken() string {
	return t.tokenRefresh
}

func (t *Token) JWTExpired() bool {
	return t.expiredAt-time.Now().Unix() <= 0
}
