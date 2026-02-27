package account

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/primary/http-adapter/codec"
	"server/internal/app/adapters/primary/http-adapter/constants"
	errorMapper "server/internal/app/adapters/primary/http-adapter/error-mapper/account_usecase"
	domain "server/internal/app/domain/account"
	"server/internal/pkg/logger"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UpdateAccountRequest struct {
	ServiceName string `json:"service_name"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	AccountId   int64  `json:"account_id"`
}

func (h *HttpHandler) UpdateAccountObj(w http.ResponseWriter, r *http.Request) {
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

	accountID, err := strconv.ParseInt(urlId, 10, 64)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid account id")
		return
	}

	userId, ok := r.Context().Value(constants.UserIDKey).(int64)
	if !ok {
		codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "user ID not found in context")
		return
	}

	account := req.toDomain()
	account.AccountId = accountID
	account.UserId = userId

	err = h.service.UpdateAccount(r.Context(), account)
	if err != nil {
		logger.Log.Error(HandlerName, zap.Error(err))

		s, m := errorMapper.Process(err)
		codec.WriteErrorJSON(w, s, m)
		return
	}

	codec.WriteJSON(w, http.StatusOK, "updated card successfully")
	return
}

func (u *UpdateAccountRequest) toDomain() *domain.Account {
	return &domain.Account{
		ServiceName: u.ServiceName,
		UserName:    u.UserName,
		Password:    u.Password,
		AccountId:   u.AccountId,
	}
}
