package app

import (
	"github.com/SerhiiKhyzhko/bookstore-oauth-go/v2/jwtsdk"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/controllers/userController"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func mapUrls(userCtrl *userController.UserController, jwtSdk *jwtsdk.JwtManager, appEnv string) {
	// public 
	router.POST("/users", userCtrl.Create) 
	router.POST("/users/login", userCtrl.Login)

	//protected
	authorized := router.Group("/")
	authorized.Use(authMiddleware(jwtSdk))
	authorized.GET("/users", userCtrl.Get)
	authorized.PUT("/users",userCtrl.Put)
	authorized.PATCH("/users",userCtrl.Patch)
	authorized.DELETE("/users",userCtrl.Delete)
	authorized.GET("/internal/users/search", userCtrl.Search)
	 
	if appEnv == "development"{
		// swagger
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}