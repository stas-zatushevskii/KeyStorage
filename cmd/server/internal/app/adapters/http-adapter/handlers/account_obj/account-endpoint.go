package account_obj

import (
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	errorMapper "server/internal/app/adapters/http-adapter/errors/account-usecase"
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
		const HandlerName = "GetAccountList"

		urlId := chi.URLParam(r, "id")

		id, err := strconv.ParseInt(urlId, 10, 64)
		if err != nil {
			codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid account id")
			return
		}

		list, err := h.service.GetAccount(r.Context(), id)
		if err != nil {
			logger.Log.Error(HandlerName, zap.Error(err))

			s, m := errorMapper.Process(err)
			codec.WriteErrorJSON(w, s, m)
			return
		}

		codec.WriteJSON(w, http.StatusOK, list)
	}
}
