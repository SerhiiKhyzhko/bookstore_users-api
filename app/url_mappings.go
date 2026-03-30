package app

import(
	"github.com/SerhiiKhyzhko/bookstore_users-api/controllers/userController"
)

func mapUrls(userCtrl *userController.UserController) {
	router.POST("/users", userCtrl.Create) 
	router.GET("/users/:users_id", userCtrl.Get)
	router.PUT("/users/:users_id",userCtrl.Put)
	router.PATCH("/users/:users_id",userCtrl.Patch)
	router.DELETE("/users/:users_id",userCtrl.Delete)
	router.GET("/internal/users/search", userCtrl.Search)
	router.POST("/users/login", userCtrl.Login) 
}