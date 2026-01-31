package token

import (
	"crypto/rand"
	"server/internal/app/config"
	"server/internal/pkg/logger"
)

func CreateRefreshToken() string {
	return string(generateRandom(config.App.GetRefreshTokenLength()))
}

func generateRandom(size int) []byte {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	return b
}
