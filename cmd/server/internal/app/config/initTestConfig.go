package config

import (
	"sync"
	"time"
)

var testCfgOnce sync.Once

func InitTestConfig() {
	testCfgOnce.Do(func() {
		App = &AppConfig{
			JWT: JWT{
				Secret:      "test-secret-very-secret",
				Issuer:      "test-issuer",
				JWTLifetime: 60, // minutes (у тебя умножается на time.Minute)
				Refresh: RefreshToken{
					Lifetime: 24 * time.Hour,
					Length:   64,
				},
			},
			Encryption: Encryption{
				AccountObjKey:  "1234567890abcdef",
				BankCardObjKey: "1234567890abcdef",
			},
		}
	})
}
