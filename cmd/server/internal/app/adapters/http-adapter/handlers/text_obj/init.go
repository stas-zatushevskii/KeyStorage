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

	router.Get("/list", handler.GetTextList())
	router.Get("/list/{id}", handler.GetTextObj())
	router.Post("/create/", handler.CreateText())
	router.Put("/update/{id}", handler.UpdateTextObj())

	return router
}
