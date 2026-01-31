package app

import (
	domain "client/internal/domain/token"
	"sync"

	"github.com/go-resty/resty/v2"
)

type Ctx struct {
	Session *Session
	HTTP    *resty.Client
	BaseURL string
}

func NewCtx(client *resty.Client) *Ctx {
	return &Ctx{
		HTTP: client,
	}
}

type Session struct {
	mu           sync.RWMutex
	token        string
	refreshToken string
}

func (c *Ctx) CreateNewSession() {
	c.Session = &Session{
		mu: sync.RWMutex{},
	}
}

func (c *Ctx) SetToken(token *domain.Token) {
	c.Session.mu.Lock()
	defer c.Session.mu.Unlock()

	c.Session.token = token.GetJWTToken()
	c.Session.refreshToken = token.GetRefreshToken()
}

func (c *Ctx) GetToken() string {
	c.Session.mu.RLock()
	defer c.Session.mu.RUnlock()
	return c.Session.token
}

func (c *Ctx) RefreshToken() string {
	c.Session.mu.RLock()
	defer c.Session.mu.RUnlock()
	return c.Session.refreshToken
}
