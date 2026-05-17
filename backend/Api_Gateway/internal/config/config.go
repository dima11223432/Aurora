package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                string              `yaml:"env" env-default:"local"`
	StoragePass        string              `yaml:"storage_pass" env-required:"true"`
	TokenTTL           time.Duration       `yaml:"token_ttl" env-required:"true"`
	GRPC               GRPCConfig          `yaml:"grpc"`
	GRPC_GatewayConfig GRPC_Gateway_Config `yaml:"grpc-gateway"`
	Auth               Auth                `yaml:"auth"`
	RedisConfig        RedisConfig         `yaml:"redis"`
	Services          ServicesConfig      `yaml:"services"`
}

type Auth struct {
	JwtSecret     string   `yaml:"jwt_secret" env-required:"true"`
	PublicMethods []string `yaml:"public_methods" env-required:"true"`
	Cors_urls     []string `yaml:"cors_urls" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	TimeOut time.Duration `yaml:"timeout"`
}
type GRPC_Gateway_Config struct {
	Port    int           `yaml:"port"`
	TimeOut time.Duration `yaml:"timeout"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	DB       int    `yaml:"db"`
}

type ServicesConfig struct {
	SSO     ServiceConfig `yaml:"sso"`
	RECS    ServiceConfig `yaml:"recs"`
}

type ServiceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	return MustLoadByPath(path)
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
