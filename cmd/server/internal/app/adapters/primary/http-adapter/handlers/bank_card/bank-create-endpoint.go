package bank_card

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/primary/http-adapter/codec"
	"server/internal/app/adapters/primary/http-adapter/constants"
	errorMapper "server/internal/app/adapters/primary/http-adapter/error-mapper/bank_card_usecase"
	domain "server/internal/app/domain/bank_card"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type CreateBankCardRequest struct {
	BankName string `json:"bank_name"`
	PID      string `json:"pid"`
}

type CreateBankCardResponse struct {
	CardID int64 `json:"card_id"`
}

func (h *HttpHandler) CreateBankCard(w http.ResponseWriter, r *http.Request) {
	const HandlerName = "CreateBankCard"

	var (
		req  = new(CreateBankCardRequest)
		resp = new(CreateBankCardResponse)
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

	bank := req.toDomain()
	bank.UserId = userId

	id, err := h.service.CreateNewBankCard(r.Context(), bank)

	if err != nil {
		logger.Log.Error(HandlerName, zap.Error(err))

		s, m := errorMapper.Process(err)
		codec.WriteErrorJSON(w, s, m)
		return
	}

	resp.CardID = id

	codec.WriteJSON(w, http.StatusOK, resp)

}

func (req CreateBankCardRequest) toDomain() *domain.BankCard {
	return &domain.BankCard{
		Bank: req.BankName,
		Pid:  req.PID,
	}
}
