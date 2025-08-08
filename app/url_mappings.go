package app

import(
	"github.com/SerhiiKhyzhko/bookstore_users-api/controllers/users"
)

func mapUrls() {
	router.POST("/users", users.Create) 
	router.GET("/users/:users_id", users.Get)
	router.PUT("/users/:users_id",users.Update)
	router.PATCH("/users/:users_id",users.Update)
	router.DELETE("/users/:users_id",users.Delete)
	router.GET("internal/users/search", users.Search)
	router.POST("/users/login", users.Login) 
}