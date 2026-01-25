package get

import (
	"client/internal/app"
	"client/pkg/http_request_sender"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Account struct {
	ServiceName string `json:"service_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

// GetAccountByID gets single text object by id
func GetAccountByID(ctx context.Context, app *app.Ctx, id int64) (*Account, error) {
	var respData Account

	url := fmt.Sprintf("http://127.0.0.1:8080/account/list/%d", id)

	response, err := http_request_sender.SendJSONRequest(
		ctx,
		http_request_sender.GET,
		http_request_sender.SendDataCmd{
			URL:    url,
			Client: app.HTTP,
			JWT:    app.GetToken(),
		},
	)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(
			"GET %s failed: status=%d body=%s",
			url,
			response.StatusCode(),
			string(response.Body()),
		)
	}

	if err := json.Unmarshal(response.Body(), &respData); err != nil {
		return nil, fmt.Errorf("json unmarshal response: %w", err)
	}

	return &respData, nil
}
