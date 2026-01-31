package text_obj

import (
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	errorMapper "server/internal/app/adapters/http-adapter/error-mapper/text_usecase"
	"server/internal/pkg/logger"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type TextResponse struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (h *HttpHandler) GetTextObj(w http.ResponseWriter, r *http.Request) {
	const HandlerName = "GetText"

	var resp = new(TextResponse)

	urlId := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(urlId, 10, 64)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid card id")
		return
	}

	card, err := h.service.GetText(r.Context(), id)
	if err != nil {
		logger.Log.Error(HandlerName, zap.Error(err))

		s, m := errorMapper.Process(err)
		codec.WriteErrorJSON(w, s, m)
		return
	}

	resp.Title = card.Title
	resp.Text = card.Text

	codec.WriteJSON(w, http.StatusOK, resp)
}
