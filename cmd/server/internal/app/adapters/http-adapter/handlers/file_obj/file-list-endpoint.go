package file_obj

import (
	"errors"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	domain "server/internal/app/domain/file_obj"

	"github.com/go-chi/chi/v5"
)

func (h *FileHandler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := parseInt64Param(chi.URLParam(r, "userId"))
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid userId")
		return
	}

	limit := parseIntQuery(r, "limit", 50)
	offset := parseIntQuery(r, "offset", 0)

	list, err := h.uc.GetFileList(r.Context(), userID, limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidUserID):
			codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid userId")
		case errors.Is(err, domain.ErrEmptyFilesList):
			codec.WriteErrorJSON(w, http.StatusBadRequest, "")
		default:
			codec.WriteErrorJSON(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	resp := make([]fileResponse, 0, len(list))
	for _, f := range list {
		resp = append(resp, *mapFile(f))
	}

	codec.WriteJSON(w, http.StatusOK, resp)
}
