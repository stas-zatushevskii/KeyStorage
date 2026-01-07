package json

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type ValidateData struct {
	R           *http.Request
	HandlerName string
	RequestData interface{}
}

type ValidateResponse struct {
	ErrMsg  json.RawMessage `json:"error"`
	ErrCode int
}

// ReadBody getting data from request body, validating json structure, unmarshalling json, validating json data.
// If OK -> errStr = "", errCode = 0
func ReadBody(data *ValidateData) (response ValidateResponse) {
	body, err := io.ReadAll(data.R.Body)
	if err != nil {
		return ValidateResponse{
			ErrMsg:  ErrorAsJSON(err),
			ErrCode: http.StatusInternalServerError,
		}
	}
	// safe Body.Close()
	defer func() {
		if err := data.R.Body.Close(); err != nil {
			logger.Log.Error(fmt.Sprintf("%s: error closing body", data.HandlerName), zap.Error(err))
		}
	}()

	if !json.Valid(body) {
		return ValidateResponse{
			ErrMsg:  json.RawMessage(`{"error": "invalid request body"}`),
			ErrCode: http.StatusBadRequest,
		}
	}

	err = json.Unmarshal(body, data.RequestData)
	if err != nil {
		return ValidateResponse{
			ErrMsg:  ErrorAsJSON(err),
			ErrCode: http.StatusInternalServerError,
		}
	}
	return
}
