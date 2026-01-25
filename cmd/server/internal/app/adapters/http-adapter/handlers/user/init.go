package user_obj

import (
	"context"
	"server/internal/app/adapters/http-adapter/middlewares"
	"server/internal/pkg/token"

	"github.com/go-chi/chi/v5"
)

type service interface {
	RegisterNewUser(ctx context.Context, username, password string) (*token.Tokens, error)
	Login(ctx context.Context, username, password string) (*token.Tokens, error)
	Authenticate(token string) (int64, error)
	RefreshJWTToken(ctx context.Context, jwt, refreshToken string) (*token.Tokens, error)
}

type httpHandler struct {
	service service
}

func newHandler(s service) *httpHandler {
	return &httpHandler{
		service: s,
	}
}

func New(service service) *chi.Mux {
	r := chi.NewRouter()

	handler := newHandler(service)

	r.Post("/auth/register", handler.RegistrationHandler())
	r.Post("/auth/login", handler.LoginHandler())
	r.With(middlewares.JWTMiddleware(service)).Post("/auth/refresh", handler.RefreshTokenHandler())

	return r
}
