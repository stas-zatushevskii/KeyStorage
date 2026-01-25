package get

import (
	"client/internal/app"
	"client/pkg/http_request_sender"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Text struct {
	ID    int64  `json:"text_id"`
	Title string `json:"title"`
}

// GetTextList gets list of text objects
func GetTextList(ctx context.Context, app *app.Ctx) ([]Text, error) {
	var respData []Text

	response, err := http_request_sender.SendJSONRequest(
		ctx,
		http_request_sender.GET,
		http_request_sender.SendDataCmd{
			URL:    "http://127.0.0.1:8080/text/list",
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
