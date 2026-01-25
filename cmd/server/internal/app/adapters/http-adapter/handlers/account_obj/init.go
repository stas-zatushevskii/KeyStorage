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

type httpHandler struct {
	service service
}

func newHandler(s service) *httpHandler {
	return &httpHandler{
		service: s,
	}
}

func New(service service) *chi.Mux {
	router := chi.NewRouter()

	handler := newHandler(service)

	router.Get("/list", handler.GetAccountList())
	router.Get("/list/{id}", handler.GetAccountObj())
	router.Post("/create", handler.CreateAccount())
	router.Put("/update/{id}", handler.UpdateAccountObj())

	return router
}
