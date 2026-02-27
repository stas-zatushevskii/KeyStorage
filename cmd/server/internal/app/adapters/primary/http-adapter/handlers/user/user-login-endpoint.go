package user_obj

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/primary/http-adapter/codec"
	errorMapper "server/internal/app/adapters/primary/http-adapter/error-mapper/user_usecase"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Error        string `json:"error"`
}

func (h *HttpHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	const HandlerName = "LoginHandler"

	var (
		req  = new(LoginRequest)
		resp = new(LoginResponse)
	)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "json decode error")
		return
	}

	tokens, err := h.service.Login(r.Context(), req.Username, req.Password)
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
