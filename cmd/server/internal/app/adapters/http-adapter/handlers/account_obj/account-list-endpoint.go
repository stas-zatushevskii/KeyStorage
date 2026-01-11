package account_obj

import (
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	"server/internal/app/adapters/http-adapter/constants"
	errorMapper "server/internal/app/adapters/http-adapter/error-mapper/account_usecase"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type GetAccountListResponse struct {
	Accounts []struct {
		AccountID   int64  `json:"account_id"`
		ServiceName string `json:"service_name"`
	} `json:"accounts"`
}

func (h *httpHandler) GetAccountList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "GetAccountList"

		userId := r.Context().Value(constants.UserIDKey).(int64)

		list, err := h.service.GetAccountsList(r.Context(), userId)
		if err != nil {
			logger.Log.Error(HandlerName, zap.Error(err))

			s, m := errorMapper.Process(err)
			codec.WriteErrorJSON(w, s, m)
			return
		}

		codec.WriteJSON(w, http.StatusOK, list)
	}
}
