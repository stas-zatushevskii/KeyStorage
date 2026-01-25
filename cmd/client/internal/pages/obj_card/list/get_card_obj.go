package get

import (
	"client/internal/app"
	"client/pkg/http_request_sender"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Card struct {
	CardID   int64  `json:"card_id"`
	BankName string `json:"bank_name"`
}

// GetCardList gets list of text objects
func GetCardList(ctx context.Context, app *app.Ctx) ([]Card, error) {
	var respData []Card

	response, err := http_request_sender.SendJSONRequest(
		ctx,
		http_request_sender.GET,
		http_request_sender.SendDataCmd{
			URL:    "http://127.0.0.1:8080/card/list",
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
