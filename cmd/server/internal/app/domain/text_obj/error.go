package text_obj

import "errors"

var (
	ErrTextNotFound   = errors.New("text not found")
	ErrEmptyTextsList = errors.New("empty texts list")

	ErrInvalidUserID = errors.New("invalid user id")
	ErrInvalidTextID = errors.New("invalid text id")
	ErrEmptyTitle    = errors.New("title is empty")
	ErrEmptyText     = errors.New("text is empty")

	ErrFailedCreateText        = errors.New("failed to create text")
	ErrFailedUpdateText        = errors.New("failed to update text")
	ErrTextInformationNotFound = errors.New("text information not found")
)
