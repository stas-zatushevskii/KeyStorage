package user

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	errorMapper "server/internal/app/adapters/http-adapter/error-mapper/user-usecase"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type RegisterNewUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterNewUserResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *httpHandler) RegistrationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "RegisterNewUser"

		var (
			req  = new(RegisterNewUserRequest)
			resp = new(RegisterNewUserResponse)
		)
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "json decode error")
			return
		}

		tokens, err := h.service.RegisterNewUser(r.Context(), req.Username, req.Password)
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
