package create

import (
	"client/internal/app"
	"client/pkg/http_request_sender"
	"context"
	"errors"
	"net/http"
)

type createAccountRequest struct {
	ServiceName string `json:"service_name"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
}

// CreateAccountObj create new User, get tokens
func CreateAccountObj(app *app.Ctx, serviceName, userName, password string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reqData := new(createAccountRequest)

	reqData.UserName = userName
	reqData.Password = password
	reqData.ServiceName = serviceName

	response, err := http_request_sender.SendJSONRequest(ctx, http_request_sender.POST, http_request_sender.SendDataCmd{
		URL:    "http://127.0.0.1:8080/account/create",
		Data:   reqData,
		Client: app.HTTP,
		JWT:    app.GetToken(),
	})

	if err != nil {
		return err // fixme: add custom err
	}

	if response.StatusCode() != http.StatusOK {
		return errors.New(string(response.Body()))
	}

	return nil
}
