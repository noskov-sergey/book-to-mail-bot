package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Telegram `yaml:"telegram"`
	Mail     `yaml:"mail"`
	Path     string `yaml:"storage"`
	Size     int    `yaml:"bath_size"`
}

type Telegram struct {
	Token string `yaml:"token"`
	Host  string `yaml:"host"`
}

type Mail struct {
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	From     string   `yaml:"from"`
	To       []string `yaml:"to"`
	Password string   `yaml:"password"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s, [ERR] %s", configPath, err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config: %s, [ERR] %s", configPath, err)
	}

	return &cfg
}
