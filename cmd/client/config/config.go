package config

var App *Config

type Config struct {
	Auth
}

type Auth struct {
	JWTRefreshToken string
	JWTAccessToken  string
}
