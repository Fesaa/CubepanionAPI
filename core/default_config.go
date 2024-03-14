package core

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DefaultConfig struct {
	YServiceName string `yaml:"service_name"`
	YHost        string `yaml:"host"`
	YPort        int    `yaml:"port"`

	YDatabase DefaultDatabaseConfig `yaml:"database"`
	YRedis    DefaultRedisConfig    `yaml:"redis"`
}

func (c DefaultConfig) ServiceName() string {
	return c.YServiceName
}

func (c DefaultConfig) Host() string {
	return c.YHost
}

func (c DefaultConfig) Port() int {
	return c.YPort
}

func (c DefaultConfig) Database() DatabaseConfig {
	return &c.YDatabase
}

func (c DefaultConfig) Redis() RedisConfig {
	return c.YRedis
}

type DefaultDatabaseConfig struct {
	YHost     string `yaml:"host"`
	YPort     int    `yaml:"port"`
	YUsername string `yaml:"username"`
	YPassword string `yaml:"password"`
	YDatabase string `yaml:"database"`
	YSslmode  string `yaml:"sslmode"`
}

func (c DefaultDatabaseConfig) Host() string {
	return c.YHost
}

func (c DefaultDatabaseConfig) Port() int {
	return c.YPort
}

func (c DefaultDatabaseConfig) Username() string {
	return c.YUsername
}

func (c DefaultDatabaseConfig) Password() string {
	return c.YPassword
}

func (c DefaultDatabaseConfig) Database() string {
	return c.YDatabase
}

func (c DefaultDatabaseConfig) SSLMode() string {
	return c.YSslmode
}

func (d DefaultDatabaseConfig) AsConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", d.Username(), d.Password(), d.Host(), d.Port(), d.Database(), d.SSLMode())
}

type DefaultRedisConfig struct {
	YHost     string `yaml:"host"`
	YPassWord string `yaml:"password"`
	YDB       int    `yaml:"db"`
}

func (c DefaultRedisConfig) Host() string {
	return c.YHost
}

func (c DefaultRedisConfig) Password() string {
	return c.YPassWord
}

func (c DefaultRedisConfig) DB() int {
	return c.YDB
}

func LoadDefaultConfigFromEnv() MicroServiceConfig {
	return &DefaultConfig{
		YServiceName: os.Getenv("SERVICE_NAME"),
		YHost:        os.Getenv("HOST"),
		YPort:        8080,
		YDatabase: DefaultDatabaseConfig{
			YHost:     os.Getenv("DB_HOST"),
			YPort:     5432,
			YUsername: os.Getenv("DB_USERNAME"),
			YPassword: os.Getenv("DB_PASSWORD"),
			YDatabase: os.Getenv("DB_NAME"),
			YSslmode:  "disable",
		},
		YRedis: DefaultRedisConfig{
			YHost:     os.Getenv("REDIS_HOST"),
			YPassWord: os.Getenv("REDIS_PASSWORD"),
			YDB:       0,
		},
	}
}

func LoadDefaultConfig(path string) (MicroServiceConfig, error) {
	config := DefaultConfig{}
	err := LoadConfig(path, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func LoadConfig[T any](path string, config *T) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return err
	}

	return nil
}
