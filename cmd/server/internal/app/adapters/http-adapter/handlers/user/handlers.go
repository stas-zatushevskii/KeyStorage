package user

import (
	"net/http"
	"server/internal/app/adapters/http-adapter/errors/user"
	"server/internal/app/adapters/http-adapter/json"
)

type HTTPHandler struct {
	service service
}

func NewHandler(s service) *HTTPHandler {
	return &HTTPHandler{
		service: s,
	}
}

func (h *HTTPHandler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "LoginHandler"

		var (
			req  = new(RegisterNewUserRequest)
			resp = new(RegisterNewUserResponse)
		)
		parseResult := json.ReadBody(&json.ValidateData{HandlerName: HandlerName, RequestData: req, R: r})

		if parseResult.ErrCode != 0 {
			json.WriteJSONResponse(w, parseResult.ErrCode, parseResult.ErrMsg)
			return
		}

		tokens, err := h.service.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			response := user.ProcessServiceErrors(err, HandlerName)
			json.WriteJSONResponse(w, response.HTTPStatus, response.ErrMsg)
			return
		}

		resp.Token = tokens.JWTToken
		resp.RefreshToken = tokens.RefreshToken

		json.WriteJSONResponse(w, http.StatusOK, resp)
	}
}

func (h *HTTPHandler) RefreshTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "RefreshTokenHandler"

		var (
			req  = new(RegisterNewUserRequest)
			resp = new(RegisterNewUserResponse)
		)
		parseResult := json.ReadBody(&json.ValidateData{HandlerName: HandlerName, RequestData: req, R: r})

		if parseResult.ErrCode != 0 {
			json.WriteJSONResponse(w, parseResult.ErrCode, parseResult.ErrMsg)
			return
		}

		tokens, err := h.service.RefreshJWTToken(r.Context(), req.Username, req.Password)
		if err != nil {
			response := user.ProcessServiceErrors(err, HandlerName)
			json.WriteJSONResponse(w, response.HTTPStatus, response.ErrMsg)
			return
		}

		resp.Token = tokens.JWTToken
		resp.RefreshToken = tokens.RefreshToken

		json.WriteJSONResponse(w, http.StatusOK, resp)
	}
}

func (h *HTTPHandler) RegistrationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "RegisterNewUser"

		var (
			req  = new(RegisterNewUserRequest)
			resp = new(RegisterNewUserResponse)
		)
		parseResult := json.ReadBody(&json.ValidateData{HandlerName: HandlerName, RequestData: req, R: r})

		if parseResult.ErrCode != 0 {
			json.WriteJSONResponse(w, parseResult.ErrCode, parseResult.ErrMsg)
			return
		}

		tokens, err := h.service.RegisterNewUser(r.Context(), req.Username, req.Password)
		if err != nil {
			response := user.ProcessServiceErrors(err, HandlerName)
			json.WriteJSONResponse(w, response.HTTPStatus, response.ErrMsg)
			return
		}

		resp.Token = tokens.JWTToken
		resp.RefreshToken = tokens.RefreshToken

		json.WriteJSONResponse(w, http.StatusOK, resp)
	}
}
