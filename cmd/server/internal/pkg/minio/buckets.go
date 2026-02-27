package minio

import (
	"context"
	"server/internal/app/config"

	"github.com/minio/minio-go/v7"
)

func (mc *MinioAdapter) InitMinio() error {
	// Создание контекста с возможностью отмены операции
	ctx := context.Background()

	// Проверка наличия бакета и его создание, если не существует
	exists, err := mc.mc.BucketExists(ctx, config.App.GetMinioBucketName())
	if err != nil {
		return err
	}
	if !exists {
		err := mc.mc.MakeBucket(ctx, config.App.GetMinioBucketName(), minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
