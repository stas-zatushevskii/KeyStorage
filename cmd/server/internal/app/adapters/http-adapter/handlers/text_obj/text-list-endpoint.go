package text_obj

import (
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	"server/internal/app/adapters/http-adapter/constants"
	errorMapper "server/internal/app/adapters/http-adapter/error-mapper/bank_card_usecase"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type Text struct {
	TextID int64  `json:"text_id"`
	Title  string `json:"title"`
}

func (h *HttpHandler) GetTextList(w http.ResponseWriter, r *http.Request) {
	const HandlerName = "GetTextList"
	resp := make([]Text, 0)

	userId := r.Context().Value(constants.UserIDKey).(int64)

	list, err := h.service.GetTextList(r.Context(), userId)
	if err != nil {
		logger.Log.Error(HandlerName, zap.Error(err))

		s, m := errorMapper.Process(err)
		codec.WriteErrorJSON(w, s, m)
		return
	}

	for _, item := range list {
		c := Text{
			Title:  item.Title,
			TextID: item.TextId,
		}
		resp = append(resp, c)
	}

	codec.WriteJSON(w, http.StatusOK, resp)

}
