package file_obj

import (
	"errors"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	domain "server/internal/app/domain/file_obj"

	"github.com/go-chi/chi/v5"
)

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt64Param(chi.URLParam(r, "id"))
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.uc.Delete(r.Context(), id)
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

	w.WriteHeader(http.StatusNoContent)
}
