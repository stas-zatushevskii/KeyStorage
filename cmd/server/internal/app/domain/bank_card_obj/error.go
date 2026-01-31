package bank_card_obj

import "errors"

var (
	ErrFaildeCreateBankCardObject = errors.New("fail create card object")

	ErrInvalidUserID = errors.New("invalid user id")
	ErrInvalidCardID = errors.New("invalid card id")
	ErrEmptyBankName = errors.New("bank name is empty")
	ErrEmptyPID      = errors.New("pid is empty")

	ErrFailedUpdateBankCard = errors.New("failed to update card object")
	ErrBankCardNotFound     = errors.New("bank card not found")
)
