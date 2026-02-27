package user_obj

import (
	"context"
	"server/internal/app/adapters/primary/http-adapter/middlewares"
	"server/internal/pkg/token"

	"github.com/go-chi/chi/v5"
)

type service interface {
	RegisterNewUser(ctx context.Context, username, password string) (*token.Tokens, error)
	Login(ctx context.Context, username, password string) (*token.Tokens, error)
	Authenticate(token string) (int64, error)
	RefreshJWTToken(ctx context.Context, jwt, refreshToken string) (*token.Tokens, error)
}

type HttpHandler struct {
	service service
}

func New(s service) *HttpHandler {
	return &HttpHandler{
		service: s,
	}
}

func (h *HttpHandler) Routes(service service) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/auth/register", h.RegistrationHandler)
	r.Post("/auth/login", h.LoginHandler)
	r.With(middlewares.JWTMiddleware(service)).Post("/auth/refresh", h.RefreshTokenHandler)

	return r
}
