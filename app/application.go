package app

import (
	"github.com/SerhiiKhyzhko/bookstore-oauth-go/v2/jwtsdk"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/controllers/userController"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication(port string, controller *userController.UserController, jwtSdk *jwtsdk.JwtManager, appEnv string, logger *logger.Logger) {
	mapUrls(controller, jwtSdk, appEnv)

	logger.Info("about to start the application...")
	router.Run(port)
}
