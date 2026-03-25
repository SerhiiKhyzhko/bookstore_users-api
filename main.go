package main

import (
	"log"

	"github.com/SerhiiKhyzhko/bookstore_users-api/app"
	"github.com/SerhiiKhyzhko/bookstore_users-api/config"
	"github.com/SerhiiKhyzhko/bookstore_users-api/controllers/userController"
	"github.com/SerhiiKhyzhko/bookstore_users-api/datasources/mysql/users_db"
	"github.com/SerhiiKhyzhko/bookstore_users-api/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/services"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load();err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	
	cfg := config.Load()
	loggerCfg := logger.Config{
		Level: cfg.Logger.Level,
		OutputPaths: []string{cfg.Logger.OutputPath},
	}

	logger, err := logger.NewLogger(loggerCfg)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	client, err  := usersdb.NewClient(cfg.Db)
	if err != nil {
		panic(err)
	}
	userDao := users.NewUserDao(client, logger)
	userService := services.NewUserService(userDao)
	userCtrl := userController.NewUserController(userService)
	app.StartApplication(userCtrl, logger)
}