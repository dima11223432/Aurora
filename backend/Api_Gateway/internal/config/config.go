// Package config provides configuration loading for the API Gateway service.
// It supports loading from YAML files and environment variables.
package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config holds all configuration for the API Gateway service.
type Config struct {
	Env                string              `yaml:"env" env:"ENV" env-default:"local"`
	StoragePass        string              `yaml:"storage_pass" env:"STORAGE_PASS" env-required:"true"`
	TokenTTL           time.Duration       `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"1h"`
	GRPC               GRPCConfig          `yaml:"grpc" env:"GRPC"`
	GRPC_GatewayConfig GRPC_Gateway_Config `yaml:"grpc-gateway" env:"GRPC_GATEWAY"`
	Auth               Auth                `yaml:"auth" env:"AUTH"`
	RedisConfig        RedisConfig         `yaml:"redis"`
	Services           ServiceConfig       `yaml:"services"`
}

// Auth holds JWT and CORS configuration.
type Auth struct {
	JwtSecret     string   `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
	PublicMethods []string `yaml:"public_methods" env:"PUBLIC_METHODS" env-required:"true"`
	Cors_urls     []string `yaml:"cors_urls" env:"CORS_URLS" env-required:"true"`
}

// GRPCConfig holds gRPC server configuration.
type GRPCConfig struct {
	Port    int           `yaml:"port" env:"GRPC_PORT" env-default:"44043"`
	TimeOut time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"10h"`
}
// GRPC_Gateway_Config holds gRPC-gateway HTTP server configuration.
type GRPC_Gateway_Config struct {
	Port    int           `yaml:"port" env:"GRPC_GATEWAY_PORT" env-default:"8081"`
	TimeOut time.Duration `yaml:"timeout" env:"GRPC_GATEWAY_TIMEOUT" env-default:"10h"`
}

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Host     string `yaml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	Port     int    `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
	DB       int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
}

// ServiceConfig holds connection details for downstream services (SSO, RECS).
type ServiceConfig struct {
	SSO  SSOConfig  `yaml:"sso" env-layout:"prefix"`
	RECS RECSConfig `yaml:"recs" env-layout:"prefix"`
}

// SSOConfig holds SSO service connection configuration.
type SSOConfig struct {
	Host string `yaml:"host" env:"SSO_HOST"`
	Port int    `yaml:"port" env:"SSO_PORT"`
}

// RECSConfig holds Recommendation Service connection configuration.
type RECSConfig struct {
	Host string `yaml:"host" env:"RECS_HOST"`
	Port int    `yaml:"port" env:"RECS_PORT"`
}

// MustLoad loads the configuration from a YAML file or environment variables.
// Panics if loading fails.
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
