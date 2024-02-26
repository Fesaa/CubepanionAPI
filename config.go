package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DatabaseUrl string      `yaml:"database_url"`
	Port        int         `yaml:"port"`
	Address     string      `yaml:"address"`
	RedisConfig RedisConfig `yaml:"redis"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	PassWord string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func NewConfig(path string) (*Config, error) {

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
