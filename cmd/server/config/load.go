package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

// fixme: NOT TESTED

// Config global single ton value, access to parameters through GET methods
var Config *Cfg

// LoadConfig loads single ton value Config
func LoadConfig(){
	config := new(Cfg)
	config.parseFlags()
	config.loadENV()

	if config.core.ConfigPath != "" {
		Config = loadConfigFromFile(config.core.ConfigPath)
	}

	Config = config
}

// parseFlags loads flags
func (cfg *Cfg) parseFlags() {
	var (
		serverAddrFlg = flag.String("a", "", "address to listen (e.g. 127.0.0.1:8080)")
		connPathFlag  = flag.String("d", "", "database connection URL (postgres://...)")
		configPath    = flag.String("c", "../config/config.yml", "config file path")
		debugModeFlg  = flag.Bool("t", false, "debug mode")
	)

	flag.Parse()

	if serverAddrFlg != nil {
		cfg.server.addr = *serverAddrFlg
	}
	if connPathFlag != nil {
		cfg.db.DSN = *connPathFlag
	}
	if configPath != nil {
		cfg.core.ConfigPath = *configPath
	}
	if debugModeFlg != nil {
		cfg.core.DebugMode = *debugModeFlg
	}

}

// loadENV loads virtual environment
func (cfg *Cfg) loadENV() {
	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.server.addr = addr
	}
	// TODO: RFC ?
	//if key, ok := os.LookupEnv("KEY"); ok {
	//	cfg.HashKey = key
	//}
	if dsn, ok := os.LookupEnv("DSN"); ok {
		cfg.db.DSN = dsn
	}
	if configPath, ok := os.LookupEnv("CONFIG_FILE"); ok {
		cfg.core.ConfigPath = configPath
	}

}

// loadConfigFromFile loads data from config.yml file
func loadConfigFromFile(configPath string) *Cfg {
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	config := new(Cfg)

	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}
	return config
}

