package account_obj

import (
	"context"
	domain "server/internal/app/domain/account_obj"

	"github.com/go-chi/chi/v5"
)

type service interface {
	GetAccountsList(ctx context.Context, userId int64) ([]*domain.Account, error)
	GetAccount(ctx context.Context, accountId int64) (*domain.Account, error)
	CreateNewAccountObj(ctx context.Context, account *domain.Account) (int64, error)
	UpdateAccount(ctx context.Context, account *domain.Account) error
}

type HttpHandler struct {
	service service
}

func New(service service) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

func (h *HttpHandler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/list", h.GetAccountList)
	router.Get("/list/{id}", h.GetAccountObj)
	router.Post("/create", h.CreateAccount)
	router.Put("/update/{id}", h.UpdateAccountObj)

	return router
}
