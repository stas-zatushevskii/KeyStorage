package create

import (
	"client/internal/app"
	"client/pkg/http_request_sender"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type createTextObjRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type createTextObjResponse struct {
	TextID int64 `json:"text_id"`
}

// CreateTextObj create new User, get tokens
func CreateTextObj(app *app.Ctx, title, text string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		reqData  = new(createTextObjRequest)
		respData = new(createTextObjResponse)
	)
	reqData.Title = title
	reqData.Text = text

	response, err := http_request_sender.SendJSONRequest(ctx, http_request_sender.POST, http_request_sender.SendDataCmd{
		URL:    "http://127.0.0.1:8080/text/create",
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

	err = json.Unmarshal(response.Body(), respData)
	if err != nil {
		return fmt.Errorf(`json unmarshal response: %s`, err)
	}

	return nil
}
