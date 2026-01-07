package token

import (
	"server/internal/app/config"
	"time"
)

type Tokens struct {
	UserId            int64
	JWTToken          string
	JWTExpAp          time.Time
	RefreshToken      string
	RefreshTokenExpAt time.Time
	Revoked           bool
}

func NewTokens(userId int64) *Tokens {
	return &Tokens{
		UserId: userId,
	}
}

func (t *Tokens) AddJWTToken(token string) {
	t.JWTToken = token
	t.JWTExpAp = time.Now().Add(config.App.GetJWTLifetime())
}

func (t *Tokens) AddRefreshToken(token string) {
	t.RefreshToken = token
	t.RefreshTokenExpAt = time.Now().Add(config.App.GetRefreshTokenLifeTime())
}
