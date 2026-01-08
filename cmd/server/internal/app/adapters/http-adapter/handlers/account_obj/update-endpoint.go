package account_obj

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	errorMapper "server/internal/app/adapters/http-adapter/errors/account-usecase"
	domain "server/internal/app/domain/account_obj"
	"server/internal/pkg/logger"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UpdateAccountRequest struct {
	ServiceName string `json:"service_name"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
}

func (h *httpHandler) UpdateAccountObj() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "UpdateAccountObj"

		var (
			req = new(UpdateAccountRequest)
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
		account.AccountId = id

		err = h.service.UpdateAccount(r.Context(), account)
		if err != nil {
			logger.Log.Error(HandlerName, zap.Error(err))

			s, m := errorMapper.Process(err)
			codec.WriteErrorJSON(w, s, m)
			return
		}

		codec.WriteJSON(w, http.StatusOK, "updated account successfully")
	}
}

func (u *UpdateAccountRequest) toDomain() *domain.Account {
	return &domain.Account{
		ServiceName: u.ServiceName,
		UserName:    u.UserName,
		Password:    u.Password,
	}
}
