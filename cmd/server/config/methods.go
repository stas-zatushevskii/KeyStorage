package config

func (cfg *Cfg) GetDSN() string {
	return cfg.db.DSN
}

func (cfg *Cfg) GetDebugMode() bool {
	return cfg.core.DebugMode
}

func (cfg *Cfg) GetServerAddr() string {
	return cfg.server.addr
}

func (cfg *Cfg) GetServerPort() int {
	return cfg.server.port
}

func (cfg *Cfg) GetServerHost() string {
	return cfg.server.host
}
