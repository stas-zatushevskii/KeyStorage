package bank_card_obj

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	errorMapper "server/internal/app/adapters/http-adapter/error-mapper/bank_card_usecase"
	domain "server/internal/app/domain/bank_card_obj"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type CreateBankCardRequest struct {
	UserID   int64  `json:"user_id"`
	BankName string `json:"bank_name"`
	PID      string `json:"pid"`
}

type CreateBankCardResponse struct {
	CardID int64 `json:"card_id"`
}

func (h *httpHandler) CreateBankCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		id, err := h.service.CreateNewBankCardObj(r.Context(), req.toDomain())

		if err != nil {
			logger.Log.Error(HandlerName, zap.Error(err))

			s, m := errorMapper.Process(err)
			codec.WriteErrorJSON(w, s, m)
			return
		}

		resp.CardID = id

		codec.WriteJSON(w, http.StatusOK, resp)
	}
}

func (req CreateBankCardRequest) toDomain() *domain.BankCard {
	return &domain.BankCard{
		UserId: req.UserID,
		Bank:   req.BankName,
		Pid:    req.PID,
	}
}
