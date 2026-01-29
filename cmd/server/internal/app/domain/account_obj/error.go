package account_obj

import "errors"

var (
	ErrAccountInformationNotFound = errors.New("account information not found")
	ErrEmptyAccountsList          = errors.New("empty accounts list")
	ErrFaildeCreateAccountObject  = errors.New("fail create account object")
	ErrAccountNotFound            = errors.New("account not found")
	ErrInvalidUserID              = errors.New("invalid user id")
	ErrInvalidAccountID           = errors.New("invalid account id")
	ErrEmptyServiceName           = errors.New("service name is empty")
	ErrFailedCreateAccount        = errors.New("failed to create account")
	ErrFailedUpdateAccount        = errors.New("failed to update account")
)
