package users

import (
	"net/http"

	"github.com/SerhiiKhyzhko/bookstore_users-api/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/services"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context){
	var user users.User
	
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.RestErr{
			Message: "invalid JSON body",
			Status: http.StatusBadRequest,
			Error: "bad_request",
		}
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveError := services.CreateUser(user)
	if saveError != nil {
		c.JSON(saveError.Status, saveError)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Need to implement")
}