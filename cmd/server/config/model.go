package config

type appConfig struct {
	core   core   `yaml:"Core"`
	db     db     `yaml:"Db"`
	logger logger `yaml:"Logger"`
	server server `yaml:"Server"`
}

type core struct {
	DebugMode  bool `yaml:"debug_mode"`
	ConfigPath string
}

type server struct {
	host string `yaml:"Host"`
	port int    `yaml:"Port"`
	addr string `yaml:"Addr"`
}

type db struct {
	Host            string `yaml:"Host"`
	PortRO          int    `yaml:"Port_RO"`
	PortRW          int    `yaml:"Port_RW"`
	User            string `yaml:"User"`
	Password        string `yaml:"Password"`
	Encode          string `yaml:"Encode"`
	Db              string `yaml:"DB"`
	MaxConnActive   int32  `yaml:"MaxConnActive"`
	MaxIdleConns    int32  `yaml:"MaxIdleConns"`
	ConnMaxLifetime int32  `yaml:"ConnMaxLifetime"`
	DSN             string `yaml:"DSN"`
}

type logger struct {
}
