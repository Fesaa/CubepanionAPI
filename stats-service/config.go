package main

import (
	"github.com/Fesaa/CubepanionAPI/core"
)

type StatsServiceConfig struct {
	YService string `yaml:"service_name"`
	YHost    string `yaml:"host"`
	YPort    int    `yaml:"port"`

	YGamesServiceURL string `yaml:"games_service"`

	YDatabase core.DefaultDatabaseConfig `yaml:"database"`
	YRedis    core.DefaultRedisConfig    `yaml:"redis"`
	YLogging  core.DefaultLoggingConfig  `yaml:"logging"`
}

func (c StatsServiceConfig) ServiceName() string {
	return c.YService
}

func (c StatsServiceConfig) Host() string {
	return c.YHost
}

func (c StatsServiceConfig) Port() int {
	return c.YPort
}

func (c StatsServiceConfig) GamesService() string {
	return c.YGamesServiceURL
}

func (c StatsServiceConfig) Database() core.DatabaseConfig {
	return c.YDatabase
}

func (c StatsServiceConfig) Redis() core.RedisConfig {
	return c.YRedis
}

func (c StatsServiceConfig) LoggingConfig() core.LoggingConfig {
	return c.YLogging
}
