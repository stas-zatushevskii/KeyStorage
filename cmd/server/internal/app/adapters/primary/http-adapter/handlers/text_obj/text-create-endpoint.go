package text_obj

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/primary/http-adapter/codec"
	"server/internal/app/adapters/primary/http-adapter/constants"
	errorMapper "server/internal/app/adapters/primary/http-adapter/error-mapper/text_usecase"
	domain "server/internal/app/domain/text_obj"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type CreateTextRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type CreateTextResponse struct {
	TextID int64 `json:"text_id"`
}

func (h *HttpHandler) CreateText(w http.ResponseWriter, r *http.Request) {
	const HandlerName = "CreateText"

	var (
		req  = new(CreateTextRequest)
		resp = new(CreateTextResponse)
	)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "json decode error")
		return
	}

	userId, ok := r.Context().Value(constants.UserIDKey).(int64)
	if !ok {
		codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "user ID not found in context")
		return
	}

	card := req.toDomain()
	card.UserId = userId

	id, err := h.service.CreateNewTextObj(r.Context(), card)

	if err != nil {
		logger.Log.Error(HandlerName, zap.Error(err))

		s, m := errorMapper.Process(err)
		codec.WriteErrorJSON(w, s, m)
		return
	}

	resp.TextID = id

	codec.WriteJSON(w, http.StatusOK, resp)

}

func (req CreateTextRequest) toDomain() *domain.Text {
	return &domain.Text{
		Title: req.Title,
		Text:  req.Text,
	}
}
