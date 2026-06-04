package main

import (
	"log"
	"os"

	"github.com/SerhiiKhyzhko/bookstore-oauth-go/v2/jwtsdk"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/app"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/config"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/controllers/userController"
	usersdb "github.com/SerhiiKhyzhko/bookstore_users-api/v2/datasources/mysql/users_db"
	_ "github.com/SerhiiKhyzhko/bookstore_users-api/v2/docs"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/services"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := config.Load()
	loggerCfg := logger.Config{
		Level:       cfg.Logger.Level,
		OutputPaths: []string{cfg.Logger.OutputPath},
	}

	logger, err := logger.NewLogger(loggerCfg)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	defer logger.Sync()

	oauthClient := jwtsdk.NewJwtManager(cfg.App.SecretKey, logger)

	client, err := usersdb.NewClient(cfg.Db)
	if err != nil {
		logger.Error(err.Error(), err)
		os.Exit(1)
	}
	userDao := users.NewUserDao(client, logger)
	userService := services.NewUserService(cfg.App.CtxTimeout, userDao)
	userCtrl := userController.NewUserController(userService)
	app.StartApplication(cfg.App.GinPort, userCtrl, oauthClient, cfg.App.AppEnv, logger)
}
