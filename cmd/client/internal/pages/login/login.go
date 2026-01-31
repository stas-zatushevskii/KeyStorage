package login

import (
	domain "client/internal/domain/token"
	"client/pkg/http_request_sender"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// Login create new User, get tokens, save in keyring
func Login(client *resty.Client, username, password string) (*domain.Token, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		reqData  = new(loginRequest)
		respData = new(loginResponse)
	)
	reqData.Username = username
	reqData.Password = password

	response, err := http_request_sender.SendJSONRequest(ctx, http_request_sender.POST, http_request_sender.SendDataCmd{
		URL:    "http://127.0.0.1:8080/user/auth/login", // fixme: create path from config
		Data:   reqData,
		Client: client,
	})

	if err != nil {
		return nil, err // fixme: add custom err
	}

	err = json.Unmarshal(response.Body(), respData)
	if err != nil {
		return nil, fmt.Errorf(`json unmarshal response: %s`, err)
	}

	tokens := domain.NewToken()

	err = tokens.SetJWTToken(respData.Token)
	if err != nil {
		return nil, err // fixme: add custom err
	}
	tokens.SetRefreshToken(respData.RefreshToken)

	return tokens, nil
}
