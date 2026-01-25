package file_obj

import (
	"strings"
	"time"
)

type FileID int64
type UserID int64

type StorageRef struct {
	BucketName string
	ObjectKey  string
}

func NewStorageRef(bucketName, objectKey string) (StorageRef, error) {
	bucketName = strings.TrimSpace(bucketName)
	objectKey = strings.TrimSpace(objectKey)

	if bucketName == "" {
		return StorageRef{}, ErrEmptyBucketName
	}
	if objectKey == "" {
		return StorageRef{}, ErrEmptyObjectKey
	}

	return StorageRef{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}, nil
}

type File struct {
	ID          FileID
	UserID      UserID
	Title       string
	Storage     StorageRef
	SizeBytes   int64
	ContentType string
	ETag        string
	CreatedAt   time.Time
}

func NewFile(
	userID UserID,
	title string,
	storage StorageRef,
	sizeBytes int64,
	contentType string,
	etag string,
) (*File, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if sizeBytes < 0 {
		return nil, ErrNegativeSizeBytes
	}

	return &File{
		ID:          0,
		UserID:      userID,
		Title:       strings.TrimSpace(title),
		Storage:     storage,
		SizeBytes:   sizeBytes,
		ContentType: strings.TrimSpace(contentType),
		ETag:        strings.TrimSpace(etag),
		CreatedAt:   time.Time{},
	}, nil
}
