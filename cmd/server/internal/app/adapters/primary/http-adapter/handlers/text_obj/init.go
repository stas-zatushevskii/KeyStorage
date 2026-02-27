package text_obj

import (
	"context"
	domain "server/internal/app/domain/text_obj"

	"github.com/go-chi/chi/v5"
)

type service interface {
	GetText(ctx context.Context, cardId int64) (*domain.Text, error)
	GetTextList(ctx context.Context, userId int64) ([]*domain.Text, error)
	CreateNewTextObj(ctx context.Context, card *domain.Text) (int64, error)
	UpdateText(ctx context.Context, card *domain.Text) error
}

type HttpHandler struct {
	service service
}

func New(uc service) *HttpHandler {
	return &HttpHandler{service: uc}
}

func (h *HttpHandler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/list", h.GetTextList)
	router.Get("/list/{id}", h.GetTextObj)
	router.Post("/create", h.CreateText)
	router.Put("/update/{id}", h.UpdateTextObj)

	return router
}
