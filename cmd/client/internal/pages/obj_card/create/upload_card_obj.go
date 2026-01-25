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

type createBankCardRequest struct {
	BankName string `json:"bank_name"`
	Pid      string `json:"pid"`
}

type createBankCardResponse struct {
	BankCardID int64 `json:"card_id"`
}

// CreateBankCardObj create new User, get tokens
func CreateBankCardObj(app *app.Ctx, bankName, pid string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		reqData  = new(createBankCardRequest)
		respData = new(createBankCardResponse)
	)
	reqData.BankName = bankName
	reqData.Pid = pid

	response, err := http_request_sender.SendJSONRequest(ctx, http_request_sender.POST, http_request_sender.SendDataCmd{
		URL:    "http://127.0.0.1:8080/card/create",
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
