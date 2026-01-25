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
	AccountID   int64  `json:"account_id"`
	ServiceName string `json:"service_name"`
	Username    string `json:"username"`
}

// GetAccountList gets list of text objects
func GetAccountList(ctx context.Context, app *app.Ctx) ([]Account, error) {
	var respData []Account

	response, err := http_request_sender.SendJSONRequest(
		ctx,
		http_request_sender.GET,
		http_request_sender.SendDataCmd{
			URL:    "http://127.0.0.1:8080/account/list",
			Client: app.HTTP,
			JWT:    app.GetToken(),
		},
	)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(
			"GET /text/list failed: status=%d body=%s",
			response.StatusCode(),
			string(response.Body()),
		)
	}

	if err := json.Unmarshal(response.Body(), &respData); err != nil {
		return nil, fmt.Errorf("json unmarshal response: %w", err)
	}

	return respData, nil
}
