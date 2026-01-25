package bank_card_obj

import (
	"context"
	domain "server/internal/app/domain/bank_card_obj"

	"github.com/go-chi/chi/v5"
)

type service interface {
	GetBankCard(ctx context.Context, cardId int64) (*domain.BankCard, error)
	GetBankCardList(ctx context.Context, userId int64) ([]*domain.BankCard, error)
	CreateNewBankCardObj(ctx context.Context, card *domain.BankCard) (int64, error)
	UpdateBankCard(ctx context.Context, card *domain.BankCard) error
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

	router.Get("/list", handler.GetBankCardList())
	router.Get("/list/{id}", handler.GetBankCardObj())
	router.Post("/create", handler.CreateBankCard())
	router.Put("/update/{id}", handler.UpdateBankCardObj())

	return router
}
