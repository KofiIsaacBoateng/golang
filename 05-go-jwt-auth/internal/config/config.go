package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MONGO_URI     string
	MONGO_DB_NAME string
	PORT          string
	JWT_SECRET    string
}

func Load() (Config, error) {
	// load config into device os
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env to OS!")
		return Config{}, fmt.Errorf("Failed to load env to OS: %w", err)
	}

	config := Config{
		MONGO_URI: os.Getenv("MONGO_URI"),
		MONGO_DB_NAME: os.Getenv("MONGO_DB_NAME"),
		PORT: os.Getenv("PORT"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
	}

	if config.MONGO_URI == "" {
		return Config{}, fmt.Errorf("Mongo uri is missing/empty in env.")
	}

	if config.MONGO_DB_NAME == "" {
		return Config{}, fmt.Errorf("Mongo db name is missing/empty in env.")
	}

	if config.PORT == "" {
		return Config{}, fmt.Errorf("Port is missing/empty in env.")
	}

	if config.JWT_SECRET == "" {
		return Config{}, fmt.Errorf("Jwt secret is missing/empty in env.")
	}

	return config, nil;
}