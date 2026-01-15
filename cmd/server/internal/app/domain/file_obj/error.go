package file_obj

import "errors"

var (
	ErrFileInformationNotFound = errors.New("file information not found")
	ErrEmptyFilesList          = errors.New("empty file list")
	ErrFailedCreateFileObject  = errors.New("fail create file object")
)
