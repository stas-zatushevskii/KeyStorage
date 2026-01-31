package config

import "time"

type AppConfig struct {
	Core       Core       `yaml:"core"`
	DB         DB         `yaml:"db"`
	Logger     Logger     `yaml:"logger"`
	Server     Server     `yaml:"server"`
	JWT        JWT        `yaml:"jwt"`
	Minio      Minio      `yaml:"minio"`
	Encryption Encryption `yaml:"encryption"`
	Uploads    Uploads    `yaml:"uploads"`
}

type Encryption struct {
	AccountObjKey  string `yaml:"account_obj_key"`
	BankCardObjKey string `yaml:"bank_card_obj_key"`
}

type Core struct {
	DebugMode       bool          `yaml:"debug_mode"`
	ConfigPath      string        `yaml:"config_path"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Addr string `yaml:"addr"`
}

type DB struct {
	Host            string        `yaml:"host"`
	PortRO          int           `yaml:"port_ro"`
	PortRW          int           `yaml:"port_rw"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Encode          string        `yaml:"encode"`
	DBName          string        `yaml:"dbname"`
	MaxConnActive   int           `yaml:"max_conn_active"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	DSN             string        `yaml:"dsn"`
}

type JWT struct {
	Secret      string        `yaml:"secret"`
	JWTLifetime time.Duration `yaml:"jwt_lifetime"`
	Issuer      string        `yaml:"issuer"`
	Refresh     RefreshToken  `yaml:"refresh"`
}

type RefreshToken struct {
	Lifetime time.Duration `yaml:"lifetime"`
	Length   int           `yaml:"length"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Minio struct {
	MinioEndpoint     string `yaml:"endpoint"`
	BucketName        string `yaml:"bucket_name"`
	MinioRootUser     string `yaml:"minio_root_user"`
	MinioRootPassword string `yaml:"minio_root_password"`
	MinioUseSSL       bool   `yaml:"minio_use_ssl"`
	MinioAccessKey    string `yaml:"accessKey"`
	MinioSecretKey    string `yaml:"secretKey"`
}

type Uploads struct {
	AllowedMimeTypes []string `yaml:"allowed_mime_types"`
}
