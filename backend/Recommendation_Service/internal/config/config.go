package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                       string        `yaml:"env" env:"ENV" env-default:"local"`
	StoragePass               string        `yaml:"storage_pass" env:"STORAGE_PASS" env-required:"true"`
	ParsingServiceStoragePass string        `yaml:"parsing_service_pass" env:"PARSING_SERVICE_PASS" env-required:"true"`
	TokenTTL                  time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"1h"`
	GRPC                      GRPCConfig    `yaml:"grpc" env:"GRPC"`
	Redis                     RedisConfig   `yaml:"redis" env:"REDIS"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env:"GRPC_PORT" env-default:"44040"`
	TimeOut time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"10h"`
}

type RedisConfig struct {
	Host     string `yaml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
	Port     int    `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		return MustLoadFromEnv()
	}
	return MustLoadByPath(path)
}

func MustLoadFromEnv() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cant read config from env: " + err.Error())
	}
	return &cfg
}

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
