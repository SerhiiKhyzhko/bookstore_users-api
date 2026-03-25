package config

import (
	"log"
	"os"
)

type Config struct {
	Db     DbConfig
	Logger loggerConfig
}

type DbConfig struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

type loggerConfig struct {
	Level      string
	OutputPath string
}

func loadDb() DbConfig {
	db := DbConfig{}
	db.DbUser = getRequiredEnv("DB_USER")
	db.DbPassword = getRequiredEnv("DB_PASSWORD")
	db.DbHost = getRequiredEnv("DB_HOST")
	db.DbPort = getRequiredEnv("DB_PORT")
	db.DbName = getRequiredEnv("DB_NAME")
	return db
}

func loadLogger() loggerConfig {
	logger := loggerConfig{}
	logger.Level = getRequiredEnv("LEVEL")
	logger.OutputPath = getRequiredEnv("OUTPUT_PATHS")
	return logger
}

func getRequiredEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("Critical environment variable %s is missing", key)
	}
	return value
}

func Load() Config {
	cfg := Config{}
	cfg.Db = loadDb()
	cfg.Logger = loadLogger()
	return cfg
}
