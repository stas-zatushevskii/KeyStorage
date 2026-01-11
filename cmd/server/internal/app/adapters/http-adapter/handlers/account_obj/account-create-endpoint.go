package account_obj

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	errorMapper "server/internal/app/adapters/http-adapter/error-mapper/account_usecase"
	domain "server/internal/app/domain/account_obj"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type CreateAccountRequest struct {
	UserID      int64  `json:"user_id"`
	ServiceName string `json:"service_name"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
}

type CreateAccountResponse struct {
	AccountID int64 `json:"account_id"`
}

func (h *httpHandler) CreateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "CreateAccount"

		var (
			req  = new(CreateAccountRequest)
			resp = new(CreateAccountResponse)
		)

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "json decode error")
			return
		}

		id, err := h.service.CreateNewAccountObj(r.Context(), req.toDomain())

		if err != nil {
			logger.Log.Error(HandlerName, zap.Error(err))

			s, m := errorMapper.Process(err)
			codec.WriteErrorJSON(w, s, m)
			return
		}

		resp.AccountID = id

		codec.WriteJSON(w, http.StatusOK, resp)
	}
}

func (req CreateAccountRequest) toDomain() *domain.Account {
	return &domain.Account{UserId: req.UserID,
		ServiceName: req.ServiceName,
		UserName:    req.UserName,
		Password:    req.Password}
}
