package account_obj

import (
	"net/http"
	"server/internal/app/adapters/primary/http-adapter/codec"
	"server/internal/app/adapters/primary/http-adapter/constants"
	errorMapper "server/internal/app/adapters/primary/http-adapter/error-mapper/account_usecase"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type Account struct {
	AccountID   int64  `json:"account_id"`
	ServiceName string `json:"service_name"`
	Username    string `json:"username"`
}

func (h *HttpHandler) GetAccountList(w http.ResponseWriter, r *http.Request) {
	const HandlerName = "GetAccountList"

	resp := make([]Account, 0)

	userId, ok := r.Context().Value(constants.UserIDKey).(int64)
	if !ok {
		codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "user ID not found in context")
		return
	}

	list, err := h.service.GetAccountsList(r.Context(), userId)
	if err != nil {
		logger.Log.Error(HandlerName, zap.Error(err))

		s, m := errorMapper.Process(err)
		codec.WriteErrorJSON(w, s, m)
		return
	}

	for _, item := range list {
		c := Account{
			AccountID:   item.AccountId,
			ServiceName: item.ServiceName,
			Username:    item.UserName,
		}
		resp = append(resp, c)
	}

	codec.WriteJSON(w, http.StatusOK, resp)
}
