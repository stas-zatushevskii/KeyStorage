package account_obj

import (
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	errorMapper "server/internal/app/adapters/http-adapter/error-mapper/account_usecase"
	"server/internal/pkg/logger"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AccountResponse struct {
	ServiceName string `json:"service_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

func (h *httpHandler) GetAccountObj() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "GetAccount"

		var resp = new(AccountResponse)

		urlId := chi.URLParam(r, "id")

		id, err := strconv.ParseInt(urlId, 10, 64)
		if err != nil {
			codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid account id")
			return
		}

		account, err := h.service.GetAccount(r.Context(), id)
		if err != nil {
			logger.Log.Error(HandlerName, zap.Error(err))

			s, m := errorMapper.Process(err)
			codec.WriteErrorJSON(w, s, m)
			return
		}

		resp.ServiceName = account.ServiceName
		resp.Username = account.UserName
		resp.Password = account.Password

		codec.WriteJSON(w, http.StatusOK, resp)
	}
}
