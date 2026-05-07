package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI string
	MongoDBName string
	ServerPort string
}

func Load()(Config, error) {

	// load config into os env.
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("Failed to load config to OS: %w", err)
	}

	mongoURI, err := extractEnv("MONGO_URI");
	if err != nil {
		return Config{}, err
	}

	mongoDBName, err := extractEnv("MONGO_DB_NAME");
	if err != nil {
		return Config{}, err
	}

	port, err := extractEnv("PORT");
	if err != nil {
		return Config{}, err
	}


	return Config{
		MongoURI: mongoURI,
		MongoDBName: mongoDBName,
		ServerPort: port,
	}, nil

}


func extractEnv(key string)(string, error) {
	cfgValue := os.Getenv(key);
	if cfgValue == "" {
		return "", fmt.Errorf("No config value found for key: %s", key)
	}

	return cfgValue, nil
}