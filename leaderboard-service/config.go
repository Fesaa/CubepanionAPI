package main

import (
	"github.com/Fesaa/CubepanionAPI/core"
)

type LeaderboardServiceConfig struct {
	YHost string `yaml:"host"`
	YPort int    `yaml:"port"`

	YGamesServiceURL string `yaml:"games_service"`

	YDatabase core.DefaultDatabaseConfig `yaml:"database"`
	YRedis    core.DefaultRedisConfig    `yaml:"redis"`
	YLogging  core.DefaultLoggingConfig  `yaml:"logging"`
}

func (c LeaderboardServiceConfig) ServiceName() string {
	return "leaderboard-service"
}

func (c LeaderboardServiceConfig) Host() string {
	return c.YHost
}

func (c LeaderboardServiceConfig) Port() int {
	return c.YPort
}

func (c LeaderboardServiceConfig) GamesService() string {
	return c.YGamesServiceURL
}

func (c LeaderboardServiceConfig) Database() core.DatabaseConfig {
	return c.YDatabase
}

func (c LeaderboardServiceConfig) Redis() core.RedisConfig {
	return c.YRedis
}

func (c LeaderboardServiceConfig) LoggingConfig() core.LoggingConfig {
	return c.YLogging
}
