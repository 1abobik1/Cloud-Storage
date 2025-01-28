package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env             string `env:"ENV" env-required:"true"`
	StoragePath     string `env:"STORAGE_PATH" env-required:"true"`
	GRPCPort        int    `env:"GRPC_PORT" env-required:"true"`
	PublicKeyPath   string `env:"PUBLIC_KEY_PATH" env-required:"true"`
	MinIoPort       string `env:"MINIO_PORT" env-required:"true"`
	AccessKeyID     string `env:"ACCESS_KEY_ID" env-required:"true"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY" env-required:"true"`
	MinIoUseSSL     bool   `env:"MINIO_USE_SSL" env-required:"true"`
}

func MustLoad() *Config {
	path := getConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exists" + path)
	}

	err := godotenv.Load(path)
	if err != nil {
		panic(fmt.Sprintf("No .env file found at %s, relying on environment variables", path))
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(fmt.Sprintf("Failed to load environment variables: %v", err))
	}

	return &cfg
}

func getConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	return res
}
