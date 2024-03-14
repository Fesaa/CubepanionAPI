package core

type MicroServiceConfig interface {
	ServiceName() string
	Host() string
	Port() int

	Database() DatabaseConfig
	Redis() RedisConfig
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

type DatabaseProvider[D any] func(config DatabaseConfig) (D, error)
