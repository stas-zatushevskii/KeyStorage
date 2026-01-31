package get

import (
	"client/internal/app"
	"client/pkg/http_request_sender"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type File struct {
	Title string `json:"title"`
	ID    int64  `json:"id"`
}

// GetFileList gets list of text objects
func GetFileList(ctx context.Context, app *app.Ctx) ([]File, error) {
	var respData []File

	response, err := http_request_sender.SendJSONRequest(
		ctx,
		http_request_sender.GET,
		http_request_sender.SendDataCmd{
			URL:    "http://127.0.0.1:8080/file/list/",
			Client: app.HTTP,
			JWT:    app.GetToken(),
		},
	)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(
			"GET /file/list: status=%d body=%s",
			response.StatusCode(),
			string(response.Body()),
		)
	}

	if err := json.Unmarshal(response.Body(), &respData); err != nil {
		return nil, fmt.Errorf("json unmarshal response: %w", err)
	}

	return respData, nil
}
