package text_obj

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	"server/internal/app/adapters/http-adapter/constants"
	errorMapper "server/internal/app/adapters/http-adapter/error-mapper/text_usecase"
	domain "server/internal/app/domain/text_obj"
	"server/internal/pkg/logger"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UpdateTextRequest struct {
	Title string `json:"title"`
	Text  string `json:"Text"`
}

func (h *HttpHandler) UpdateTextObj(w http.ResponseWriter, r *http.Request) {
	const HandlerName = "UpdateTextObj"

	var (
		req = new(UpdateTextRequest)
	)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "json decode error")
		return
	}

	userID := r.Context().Value(constants.UserIDKey).(int64)

	textID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid text id")
		return
	}

	text := req.toDomain()
	text.UserId = userID
	text.TextId = textID

	err = h.service.UpdateText(r.Context(), text)
	if err != nil {
		logger.Log.Error(HandlerName, zap.Error(err))

		s, m := errorMapper.Process(err)
		codec.WriteErrorJSON(w, s, m)
		return
	}

	codec.WriteJSON(w, http.StatusOK, "updated text successfully")

}

func (u *UpdateTextRequest) toDomain() *domain.Text {
	return &domain.Text{
		Title: u.Title,
		Text:  u.Text,
	}
}
