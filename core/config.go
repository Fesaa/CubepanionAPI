package core

import "log/slog"

type MicroServiceConfig interface {
	ServiceName() string
	Host() string
	Port() int

	Database() DatabaseConfig
	Redis() RedisConfig
	LoggingConfig() LoggingConfig
}

type LoggingConfig interface {
	LogLevel() slog.Level
	Source() bool
	Handler() string
	LogHttp() bool
}

type DatabaseConfig interface {
	Host() string
	Port() int
	Username() string
	Password() string
	Database() string
	SSLMode() string

	AsConnectionString() string
}

type RedisConfig interface {
	Host() string
	DB() int
	Password() string
}

type DatabaseProvider[D Database] func(config DatabaseConfig) (D, error)
