package account_obj

import "errors"

var (
	ErrAccountInformationNotFound = errors.New("account information not found")
	ErrEmptyAccountsList          = errors.New("empty accounts list")
)
