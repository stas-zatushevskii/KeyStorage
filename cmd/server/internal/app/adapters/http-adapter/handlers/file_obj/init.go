package file_obj

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	usecase "server/internal/app/usecases/file_obj"
)

type FileHandler struct {
	uc *usecase.UseCase
}

func NewFileHandler(uc *usecase.UseCase) *FileHandler {
	return &FileHandler{uc: uc}
}

func (h *FileHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/files/{id}", h.GetByID)
	r.Get("/users/{userId}/files", h.ListByUserID)
	r.Post("/files", h.Create)
	r.Delete("/files/{id}", h.Delete)

	return r
}

type fileResponse struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title,omitempty"`
	BucketName  string    `json:"bucket_name"`
	ObjectKey   string    `json:"object_key"`
	SizeBytes   int64     `json:"size_bytes"`
	ContentType string    `json:"content_type,omitempty"`
	ETag        string    `json:"etag,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func parseInt64Param(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func parseIntQuery(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
