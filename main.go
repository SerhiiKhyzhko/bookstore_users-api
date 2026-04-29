package main

import (
	"log"
	"os"

	"github.com/SerhiiKhyzhko/bookstore_users-api/app"
	"github.com/SerhiiKhyzhko/bookstore_users-api/config"
	"github.com/SerhiiKhyzhko/bookstore_users-api/controllers/userController"
	"github.com/SerhiiKhyzhko/bookstore_users-api/datasources/mysql/users_db"
	"github.com/SerhiiKhyzhko/bookstore_users-api/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/services"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/SerhiiKhyzhko/bookstore-oauth-go/oauth"
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

	oauthClient := oauth.NewOAuthClient(cfg.App.OauthApiBaseUrl, cfg.App.OauthTimeout)

	client, err := usersdb.NewClient(cfg.Db)
	if err != nil {
		logger.Error(err.Error(), err)
		os.Exit(1)
	}
	userDao := users.NewUserDao(client, logger)
	userService := services.NewUserService(cfg.App.CtxTimeout, userDao)
	userCtrl := userController.NewUserController(userService, oauthClient)
	app.StartApplication(cfg.App.GinPort, userCtrl, logger)
}
