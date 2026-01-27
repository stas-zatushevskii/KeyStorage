package bank_card_obj

import (
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	errorMapper "server/internal/app/adapters/http-adapter/error-mapper/bank_card_usecase"
	"server/internal/pkg/logger"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type BankCardResponse struct {
	BankName string `json:"bank_name"`
	PID      string `json:"pid"`
}

func (h *HttpHandler) GetBankCardObj(w http.ResponseWriter, r *http.Request) {
	const HandlerName = "GetBankCard"

	var resp = new(BankCardResponse)

	urlId := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(urlId, 10, 64)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid card id")
		return
	}

	card, err := h.service.GetBankCard(r.Context(), id)
	if err != nil {
		logger.Log.Error(HandlerName, zap.Error(err))

		s, m := errorMapper.Process(err)
		codec.WriteErrorJSON(w, s, m)
		return
	}

	resp.PID = card.Pid
	resp.BankName = card.Bank

	codec.WriteJSON(w, http.StatusOK, resp)
}
