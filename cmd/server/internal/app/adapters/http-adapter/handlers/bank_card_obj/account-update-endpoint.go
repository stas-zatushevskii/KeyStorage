package bank_card_obj

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	errorMapper "server/internal/app/adapters/http-adapter/error-mapper/bank_card_usecase"
	domain "server/internal/app/domain/bank_card_obj"
	"server/internal/pkg/logger"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UpdateBankCardRequest struct {
	Pid      string `json:"pid"`
	BankName string `json:"bank_name"`
}

func (h *httpHandler) UpdateBankCardObj() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "UpdateBankCardObj"

		var (
			req = new(UpdateBankCardRequest)
		)

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "json decode error")
			return
		}

		urlId := chi.URLParam(r, "id")

		id, err := strconv.ParseInt(urlId, 10, 64)
		if err != nil {
			codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid account id")
			return
		}

		account := req.toDomain()
		account.UserId = id

		err = h.service.UpdateBankCard(r.Context(), account)
		if err != nil {
			logger.Log.Error(HandlerName, zap.Error(err))

			s, m := errorMapper.Process(err)
			codec.WriteErrorJSON(w, s, m)
			return
		}

		codec.WriteJSON(w, http.StatusOK, "updated account successfully")
	}
}

func (u *UpdateBankCardRequest) toDomain() *domain.BankCard {
	return &domain.BankCard{
		Pid:  u.Pid,
		Bank: u.BankName,
	}
}
