package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Db     DbConfig
	Logger loggerConfig
	App    appConfig
}

type DbConfig struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

type appConfig struct {
	GinPort         string
	CtxTimeout      time.Duration
	SecretKey       string
	AppEnv 			string
}

type loggerConfig struct {
	Level      string
	OutputPath string
}

func loadApp() appConfig {
	app := appConfig{}
	app.GinPort = getRequiredEnv("GIN_PORT")
	app.CtxTimeout = getTimeWithDefault("CTX_TIMEOUT", "2s")
	app.SecretKey = getRequiredEnv("SECRET_KEY")
	app.AppEnv = getRequiredEnv("APP_ENV")
	return app
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

func getTimeWithDefault(key string, defaultValue string) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		res, err := time.ParseDuration(defaultValue)
		if err != nil {
			log.Fatalf("convertation of default value failed: %v", err)
		}
		return res
	}
	result, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("Invalid value for %s", key)
	}
	return result
}

func Load() Config {
	cfg := Config{}
	cfg.Db = loadDb()
	cfg.Logger = loadLogger()
	cfg.App = loadApp()
	return cfg
}
