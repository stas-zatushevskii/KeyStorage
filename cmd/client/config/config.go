package config

var App *Config

type Config struct {
	KeyRingNames KeyRingNames
}

type KeyRingNames struct {
	JWTRefreshToken string
	JWTAccessToken  string
	JWTExpiresAt    string
	UserName        string
}
