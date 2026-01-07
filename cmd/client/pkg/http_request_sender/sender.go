package http_request_sender

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
)

type (
	SendDataCmd struct {
		Client *resty.Client
		URL    string
		Data   any
		JWT    *jwt.Token
	}
	Method int
)

const (
	POST Method = iota
	GET
)

func SendJSONRequest(c context.Context, method Method, cmd SendDataCmd) (*resty.Response, error) {

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	data, err := json.Marshal(cmd.Data)
	if err != nil {
		return nil, err
	}

	req := cmd.Client.R().
		SetHeader("Content-Type", "application/json").
		SetContext(ctx).
		SetBody(data)

	if cmd.JWT != nil {
		cmd.Client.SetHeader("Authorization", cmd.JWT.Raw)
	}

	switch method {
	case POST:
		return req.Post(cmd.URL)
	case GET:
		return req.Get(cmd.URL)
	default:
		return nil, errors.New("invalid method")
	}

}

func SendFormDataRequest(c context.Context, cmd SendDataCmd) (*resty.Response, error) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	req := cmd.Client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetContext(ctx).
		SetBody(cmd.Data)

	if cmd.JWT != nil {
		cmd.Client.SetHeader("Authorization", cmd.JWT.Raw)
	}

	return req.Post(cmd.URL)
}
