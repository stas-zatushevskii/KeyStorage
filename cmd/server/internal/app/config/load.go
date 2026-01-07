package config

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// App global singleton value
var (
	App  *AppConfig
	once sync.Once
)

// GetConfig loads singleton config (file -> env -> flags)
func GetConfig() *AppConfig {
	once.Do(func() {
		cfg := &AppConfig{}

		flags := parseFlags()

		if flags.configPath != "" {
			if fileCfg, ok := loadConfigFromFile(flags.configPath); ok {
				cfg = fileCfg
				cfg.Core.ConfigPath = flags.configPath
			} else {
				cfg.Core.ConfigPath = flags.configPath
			}
		}

		applyENV(cfg)

		applyFlags(cfg, flags)

		App = cfg
	})

	return App
}

type parsedFlags struct {
	serverAddr string
	dsn        string
	configPath string
	debugMode  bool
	debugSet   bool
}

func parseFlags() parsedFlags {
	var out parsedFlags

	serverAddrFlg := flag.String("a", "", "address to listen (e.g. 127.0.0.1:8080)")
	connPathFlag := flag.String("d", "", "database dsn (e.g. postgres://...)")
	configPath := flag.String("c", "../../config/config-server.yml", "config file path")

	debugModeFlg := flag.Bool("t", false, "debug mode")

	flag.Parse()

	out.serverAddr = *serverAddrFlg
	out.dsn = *connPathFlag
	out.configPath = *configPath

	if flag.Lookup("t") != nil {
		out.debugSet = containsArg(os.Args, "-t") || containsArg(os.Args, "--t")
		out.debugMode = *debugModeFlg
	}

	return out
}

func applyFlags(cfg *AppConfig, f parsedFlags) {
	if f.serverAddr != "" {
		cfg.Server.Addr = f.serverAddr
	}
	if f.dsn != "" {
		cfg.DB.DSN = f.dsn
	}
	if f.configPath != "" {
		cfg.Core.ConfigPath = f.configPath
	}
	if f.debugSet {
		cfg.Core.DebugMode = f.debugMode
	}
}

func containsArg(args []string, target string) bool {
	for _, a := range args {
		if a == target {
			return true
		}
	}
	return false
}

// ---- env ----

func applyENV(cfg *AppConfig) {
	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.Server.Addr = addr
	}
	if dsn, ok := os.LookupEnv("DSN"); ok {
		cfg.DB.DSN = dsn
	}
	if configPath, ok := os.LookupEnv("CONFIG_FILE"); ok {
		cfg.Core.ConfigPath = configPath
	}
	if dbg, ok := os.LookupEnv("DEBUG_MODE"); ok {
		if dbg == "1" || dbg == "true" || dbg == "TRUE" {
			cfg.Core.DebugMode = true
		}
		if dbg == "0" || dbg == "false" || dbg == "FALSE" {
			cfg.Core.DebugMode = false
		}
	}
}

// loadConfigFromFile loads data from YAML file.
// returns (cfg, true) if loaded; (nil, false) if file doesn't exist.
func loadConfigFromFile(configPath string) (*AppConfig, bool) {
	fmt.Println(configPath)
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false
		}
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	cfg := &AppConfig{}
	if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
		panic(err)
	}
	return cfg, true
}
