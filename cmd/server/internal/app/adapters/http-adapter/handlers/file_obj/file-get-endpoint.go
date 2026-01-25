package file_obj

import (
	"errors"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	domain "server/internal/app/domain/file_obj"

	"github.com/go-chi/chi/v5"
)

func mapFile(f *domain.File) *fileResponse {
	return &fileResponse{
		ID:          int64(f.ID),
		UserID:      int64(f.UserID),
		Title:       f.Title,
		BucketName:  f.Storage.BucketName,
		ObjectKey:   f.Storage.ObjectKey,
		SizeBytes:   f.SizeBytes,
		ContentType: f.ContentType,
		ETag:        f.ETag,
		CreatedAt:   f.CreatedAt,
	}
}

func (h *FileHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt64Param(chi.URLParam(r, "id"))
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid id")
		return
	}

	file, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			codec.WriteErrorJSON(w, http.StatusBadRequest, "file not found")
		case errors.Is(err, domain.ErrInvalidFileID):
			codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid id")
		default:
			codec.WriteErrorJSON(w, http.StatusBadRequest, "internal error")
		}
		return
	}

	codec.WriteJSON(w, http.StatusOK, mapFile(file))
}
