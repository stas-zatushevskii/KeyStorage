package file_obj

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	usecase "server/internal/app/usecases/file_obj"
)

type FileHandler struct {
	uc *usecase.FileObj
}

func New(uc *usecase.FileObj) *FileHandler {
	return &FileHandler{uc: uc}
}

func (h *FileHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/download/{id}", h.DownloadByID)
	r.Get("/list/", h.ListByUserID)
	r.Post("/upload", h.Create)

	return r
}
