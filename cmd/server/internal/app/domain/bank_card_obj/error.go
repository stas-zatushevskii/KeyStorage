package bank_card_obj

import "errors"

var (
	ErrBankCardInformationNotFound = errors.New("card information not found")
	ErrEmptyBankCardsList          = errors.New("empty cards list")
	ErrFaildeCreateBankCardObject  = errors.New("fail create card object")
)
