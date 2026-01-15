package text_obj

import "errors"

var (
	ErrTextInformationNotFound = errors.New("text information not found")
	ErrEmptyTextsList          = errors.New("empty text list")
	ErrFailedCreateTextObject  = errors.New("fail create text object")
)
