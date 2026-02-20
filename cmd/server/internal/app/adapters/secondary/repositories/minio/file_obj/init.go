package file_obj

import (
	"github.com/minio/minio-go/v7"
)

type Repository struct {
	mc *minio.Client
}

func New(mc *minio.Client) *Repository {
	return &Repository{mc: mc}
}
