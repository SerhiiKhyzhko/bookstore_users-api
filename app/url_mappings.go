package app

import(
	"github.com/SerhiiKhyzhko/bookstore_users-api/controllers/users"
)

func mapUrls() {
	router.POST("/users", users.CreateUser) 
	router.GET("/users/:users_id", users.GetUser)
}