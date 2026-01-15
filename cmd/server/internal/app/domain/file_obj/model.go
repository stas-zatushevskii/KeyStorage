package file_obj

import "github.com/google/uuid"

type File struct {
	FileName string
	FileID   uuid.UUID // generate by module
	UserId   int64
}
