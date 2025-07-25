package app

import(
	"github.com/SerhiiKhyzhko/bookstore_users-api/controllers/users"
)

func mapUrls() {
	router.POST("/users", users.CreateUser) 
	router.GET("/users/:users_id", users.GetUser)
	router.PUT("/users/:users_id",users.UpdateUser)
	router.PATCH("/users/:users_id",users.UpdateUser)
	router.DELETE("/users/:users_id",users.DeleteUser)
}