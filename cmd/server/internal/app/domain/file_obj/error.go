package file_obj

import "errors"

var (
	ErrInvalidFileID     = errors.New("invalid file id")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrEmptyBucketName   = errors.New("empty bucket name")
	ErrEmptyObjectKey    = errors.New("empty object key")
	ErrNegativeSizeBytes = errors.New("size_bytes must be >= 0")
	ErrFileNotFound      = errors.New("file not found")
	ErrEmptyFilesList    = errors.New("empty files list")
)
