module server

go 1.24.4

require (
	github.com/rs/zerolog v1.34.0
	github.com/simukti/sqldb-logger v0.0.0-20230108155151-646c1a075551
	github.com/simukti/sqldb-logger/logadapter/zerologadapter v0.0.0-20230108155151-646c1a075551
	go.uber.org/zap v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
