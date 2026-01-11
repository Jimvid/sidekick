package config

import (
	"fmt"
	"os"
)

type Config struct {
	TABLE_NAME   string
	CLERK_SECRET string
}

var AppConfig *Config

func NewConfig() *Config {
	return &Config{
		TABLE_NAME:   MustGetEnv("TABLE_NAME"),
		CLERK_SECRET: MustGetEnv("CLERK_SECRET"),
	}
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s is required", key))
	}
	return value
}
