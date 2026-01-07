package minio

import (
	"context"
	"server/internal/app/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client  реализация интерфейса MinioClient
type Client struct {
	mc *minio.Client // Клиент Minio
}

// NewClient  создает новый экземпляр Minio Client
func NewClient() *Client {
	return &Client{} // Возвращает новый экземпляр minioClient с указанным именем бакета
}

// InitMinio подключается к Minio и создает бакет, если не существует
// Бакет - это контейнер для хранения объектов в Minio. Он представляет собой пространство имен, в котором можно хранить и организовывать файлы и папки.
func (m *Client) InitMinio() error {
	// Создание контекста с возможностью отмены операции
	ctx := context.Background()

	// Подключение к Minio с использованием имени пользователя и пароля
	client, err := minio.New(config.App.GetMinioEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(config.App.GetMinioRootUser(), config.App.GetMinioRootPassword(), ""),
		Secure: config.App.GetMinioUseSSL(),
	})
	if err != nil {
		return err
	}

	// Установка подключения Minio
	m.mc = client

	// Проверка наличия бакета и его создание, если не существует
	exists, err := m.mc.BucketExists(ctx, config.App.GetMinioBucketName())
	if err != nil {
		return err
	}
	if !exists {
		err := m.mc.MakeBucket(ctx, config.App.GetMinioBucketName(), minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
