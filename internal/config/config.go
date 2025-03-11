package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	StorageRoot string `yaml:"storage_root" env-required:"true"`
	Rules       string `yaml:"rules" env-required:"true"`
	Chart       string `yaml:"chart" env-required:"true"`
	ServerPort  string `yaml:"server_port" env-required:"true"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file %s does not exist", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read configuration file %s: %v", configPath, err)
	}

	return &cfg
}
