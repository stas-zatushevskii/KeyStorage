package bank_card_obj

import (
	"net/http"
	"server/internal/app/adapters/primary/http-adapter/codec"
	"server/internal/app/adapters/primary/http-adapter/constants"
	errorMapper "server/internal/app/adapters/primary/http-adapter/error-mapper/bank_card_usecase"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type Card struct {
	CardID   int64  `json:"card_id"`
	BankName string `json:"bank_name"`
}

func (h *HttpHandler) GetBankCardList(w http.ResponseWriter, r *http.Request) {
	const HandlerName = "GetBankCardList"
	resp := make([]Card, 0)

	userId, ok := r.Context().Value(constants.UserIDKey).(int64)
	if !ok {
		codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "user ID not found in context")
		return
	}

	list, err := h.service.GetBankCardList(r.Context(), userId)
	if err != nil {
		logger.Log.Error(HandlerName, zap.Error(err))

		s, m := errorMapper.Process(err)
		codec.WriteErrorJSON(w, s, m)
		return
	}

	for _, item := range list {
		c := Card{
			CardID:   item.CardId,
			BankName: item.Bank,
		}
		resp = append(resp, c)
	}

	codec.WriteJSON(w, http.StatusOK, resp)

}
