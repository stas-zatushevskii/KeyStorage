package file_obj

import (
	"errors"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	"server/internal/app/adapters/http-adapter/constants"
	domain "server/internal/app/domain/file_obj"
)

type fileResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

func (h *FileHandler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constants.UserIDKey).(int64)

	list, err := h.uc.GetFileList(r.Context(), userId)
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
		resp = append(
			resp,
			fileResponse{
				ID:    f.ID,
				Title: f.Title,
			})
	}

	codec.WriteJSON(w, http.StatusOK, resp)
}
