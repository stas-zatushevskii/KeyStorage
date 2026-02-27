package file

import (
	"context"
	"io"
	"net/http"
	domain "server/internal/app/domain/file"

	"github.com/go-chi/chi/v5"
)

type Service interface {
	GetFileStream(ctx context.Context, userID, fileID int64) (*domain.File, io.ReadCloser, error)
	GetFileList(ctx context.Context, userID int64) ([]*domain.File, error)
	UploadAndCreate(ctx context.Context, file *domain.File, data []byte) (int64, error)
	GetByID(ctx context.Context, fileID int64) (*domain.File, error)
}

type FileHandler struct {
	uc Service
}

func New(uc Service) *FileHandler {
	return &FileHandler{uc: uc}
}

func (h *FileHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/download/{id}", h.DownloadByID)
	r.Get("/list/", h.ListByUserID)
	r.Post("/upload", h.Create)

	return r
}
