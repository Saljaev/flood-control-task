package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type (
	Config struct {
		HTTPServer   `yaml:"http_server"`
		Redis        `yaml:"redis"`
		RequestCount int64         `yaml:"request_count"`
		Interval     time.Duration `yaml:"interval"`
	}

	HTTPServer struct {
		Address     string        `yaml:"address"`
		Timeout     time.Duration `yaml:"timeout"`
		IdleTimeout time.Duration `yaml:"idle_timeout"`
	}

	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	}
)

func ConfgiLoad() *Config {
	yamlFile, err := os.ReadFile("config/config.yaml")

	if err != nil {
		log.Fatal("failed to read config.yaml")
	}

	var cfg Config

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("unmarshal failed with error: %w", err)
	}

	return &cfg
}
