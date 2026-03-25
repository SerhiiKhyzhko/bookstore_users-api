package app

import (
	"github.com/SerhiiKhyzhko/bookstore_users-api/controllers/userController"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication(controller *userController.UserController, logger *logger.Logger) {
	mapUrls(controller)

	logger.Info("about to start the application...")
	router.Run(":8081")
}