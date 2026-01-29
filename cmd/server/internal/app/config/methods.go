package config

import "time"

// ---- DB ----

func (cfg *AppConfig) GetDSN() string {
	return cfg.DB.DSN
}

func (cfg *AppConfig) GetMaxIdleConns() int {
	return cfg.DB.MaxIdleConns
}

func (cfg *AppConfig) GetMaxOpenConns() int {
	return cfg.DB.MaxConnActive
}

func (cfg *AppConfig) GetConnMaxLifetime() time.Duration {
	return cfg.DB.ConnMaxLifetime
}

// ---- CORE ----

func (cfg *AppConfig) GetDebugMode() bool {
	return cfg.Core.DebugMode
}
func (cfg *AppConfig) GetShutDownTimeout() time.Duration {
	return cfg.Core.ShutdownTimeout
}

// ---- SERVER ----

func (cfg *AppConfig) GetServerAddr() string {
	return cfg.Server.Addr
}

func (cfg *AppConfig) GetServerPort() int {
	return cfg.Server.Port
}

func (cfg *AppConfig) GetServerHost() string {
	return cfg.Server.Host
}

// ---- JWT ----

func (cfg *AppConfig) GetJWTSecret() string {
	return cfg.JWT.Secret
}

func (cfg *AppConfig) GetIssuer() string {
	return cfg.JWT.Issuer
}

func (cfg *AppConfig) GetJWTLifetime() time.Duration {
	return cfg.JWT.JWTLifetime
}

func (cfg *AppConfig) GetRefreshTokenLifeTime() time.Duration {
	return cfg.JWT.Refresh.Lifetime
}

func (cfg *AppConfig) GetRefreshTokenLength() int {
	return cfg.JWT.Refresh.Length
}

// ---- MiniO ----

func (cfg *AppConfig) GetMinioEndpoint() string {
	return cfg.Minio.MinioEndpoint
}

func (cfg *AppConfig) GetMinioBucketName() string {
	return cfg.Minio.BucketName
}

func (cfg *AppConfig) GetMinioRootUser() string {
	return cfg.Minio.MinioRootUser
}

func (cfg *AppConfig) GetMinioRootPassword() string {
	return cfg.Minio.MinioRootPassword
}

func (cfg *AppConfig) GetMinioUseSSL() bool {
	return cfg.Minio.MinioUseSSL
}

func (cfg *AppConfig) GetMinioAccessKey() string {
	return cfg.Minio.MinioAccessKey
}

func (cfg *AppConfig) GetMinioSecretKey() string {
	return cfg.Minio.MinioSecretKey
}

// ---- Encryption ----

func (cfg *AppConfig) GetAccountObjEncryptionKey() string {
	return cfg.Encryption.AccountObjKey
}
func (cfg *AppConfig) GetBankCardObjEncryptionKey() string {
	return cfg.Encryption.BankCardObjKey
}

// ---- File Types

func (cfg *AppConfig) AllowedMimeSet() map[string]struct{} {
	m := make(map[string]struct{}, len(cfg.Uploads.AllowedMimeTypes))
	for _, t := range cfg.Uploads.AllowedMimeTypes {
		m[t] = struct{}{}
	}
	return m
}
