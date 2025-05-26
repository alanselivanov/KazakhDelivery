package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}

type MongoDBConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
	Timeout  int    `yaml:"timeout"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	TTL      int    `yaml:"ttl"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	FromName string `yaml:"from_name"`
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	MongoDB MongoDBConfig `yaml:"mongodb"`
	Redis   RedisConfig   `yaml:"redis"`
	SMTP    SMTPConfig    `yaml:"smtp"`
}

func LoadConfig() *Config {
	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(filename), "../..")
	envPath := filepath.Join(basePath, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: .env file not found or cannot be loaded: %v", err)
		log.Printf("Looking for .env at: %s", envPath)
	} else {
		log.Printf("Environment variables loaded from: %s", envPath)
	}

	return &Config{
		Server: ServerConfig{
			Port: "50053",
		},
		MongoDB: MongoDBConfig{
			URI:      "mongodb://localhost:27017",
			Database: "user_service",
			Timeout:  10,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
			TTL:      300,
		},
		SMTP: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			FromName: os.Getenv("SMTP_FROM_NAME"),
		},
	}
}
