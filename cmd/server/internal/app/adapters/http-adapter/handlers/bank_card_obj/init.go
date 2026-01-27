package bank_card_obj

import (
	"context"
	domain "server/internal/app/domain/bank_card_obj"
	usecase "server/internal/app/usecases/file_obj"

	"github.com/go-chi/chi/v5"
)

type service interface {
	GetBankCard(ctx context.Context, cardId int64) (*domain.BankCard, error)
	GetBankCardList(ctx context.Context, userId int64) ([]*domain.BankCard, error)
	CreateNewBankCardObj(ctx context.Context, card *domain.BankCard) (int64, error)
	UpdateBankCard(ctx context.Context, card *domain.BankCard) error
}

type HttpHandler struct {
	service service
}

func New(s service) *HttpHandler {
	return &HttpHandler{service: s}
}

func (h *HttpHandler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/list", h.GetBankCardList)
	router.Get("/list/{id}", h.GetBankCardObj)
	router.Post("/create", h.CreateBankCard)
	router.Put("/update/{id}", h.UpdateBankCardObj)

	return router
}
