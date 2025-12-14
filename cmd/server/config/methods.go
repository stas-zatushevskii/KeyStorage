package config

import "time"

func (cfg *appConfig) GetDSN() string {
	return cfg.db.DSN
}

func (cfg *appConfig) GetDebugMode() bool {
	return cfg.core.DebugMode
}

func (cfg *appConfig) GetServerAddr() string {
	return cfg.server.addr
}

func (cfg *appConfig) GetServerPort() int {
	return cfg.server.port
}

func (cfg *appConfig) GetServerHost() string {
	return cfg.server.host
}

func (cfg *appConfig) GetMaxIdleConns() int {
	//TODO implement me
	panic("implement me")
}

func (cfg *appConfig) GetMaxOpenConns() int {
	//TODO implement me
	panic("implement me")
}

func (cfg *appConfig) GetConnMaxLifetime() time.Duration {
	//TODO implement me
	panic("implement me")
}
