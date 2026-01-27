package upload

import (
	"client/internal/app"
	"client/pkg/http_request_sender"
	"errors"
	"net/http"
)

// UploadFileObj create new User, get tokens
func UploadFileObj(app *app.Ctx, path string) error {
	response, err := http_request_sender.SendFormDataRequest(http_request_sender.SendFileCmd{
		Client:   app.HTTP,
		FilePath: path,
		JWT:      app.GetToken(),
		URL:      "http://127.0.0.1:8080/file/upload",
	})

	if err != nil {
		return err // fixme: add custom err
	}

	if response.StatusCode() != http.StatusOK {
		return errors.New(string(response.Body()))
	}

	return nil
}
