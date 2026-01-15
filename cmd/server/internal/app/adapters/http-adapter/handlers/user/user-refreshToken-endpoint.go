package user_obj

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	errorMapper "server/internal/app/adapters/http-adapter/error-mapper/user_usecase"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type RefreshTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *httpHandler) RefreshTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "RefreshTokenHandler"

		var (
			req  = new(RefreshTokenRequest)
			resp = new(RefreshTokenResponse)
		)

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "json decode error")
			return
		}

		tokens, err := h.service.RefreshJWTToken(r.Context(), req.Username, req.Password)
		if err != nil {
			logger.Log.Error(HandlerName, zap.Error(err))

			s, m := errorMapper.Process(err)
			codec.WriteErrorJSON(w, s, m)
			return
		}

		resp.Token = tokens.JWTToken
		resp.RefreshToken = tokens.RefreshToken

		codec.WriteJSON(w, http.StatusOK, resp)
	}
}
