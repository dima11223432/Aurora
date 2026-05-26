// Package config provides configuration loading for the SSO service.
package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config holds the SSO service configuration.
type Config struct {
	Env             string        `yaml:"env" env:"ENV" env-default:"local"`
	StoragePass     string        `yaml:"storage_pass" env:"STORAGE_PASS" env-required:"true"`
	TestPostgresDNS string        `yaml:"test_postgres_dns" env:"TEST_POSTGRES_DNS" env-required:"true"`
	TokenTTL        time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"1h"`
	GRPC            GRPCConfig    `yaml:"grpc" env:"GRPC"`
}

// GRPCConfig holds gRPC server configuration.
type GRPCConfig struct {
	Port    int           `yaml:"port" env:"GRPC_PORT" env-default:"44044"`
	TimeOut time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"10h"`
}

// MustLoad loads configuration from file or environment. Panics on failure.
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		return MustLoadFromEnv()
	}
	return MustLoadByPath(path)
}

// MustLoadFromEnv loads configuration from environment variables. Panics on failure.
func MustLoadFromEnv() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cant read config from env: " + err.Error())
	}
	return &cfg
}

// MustLoadByPath loads configuration from a YAML file. Panics on failure.
func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config does not exists" + configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cant read config" + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path-to-config")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
