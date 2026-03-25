package userController

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SerhiiKhyzhko/bookstore-oauth-go/oauth"
	"github.com/SerhiiKhyzhko/bookstore_users-api/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/services"
	"github.com/SerhiiKhyzhko/bookstore_users-api/user_errors"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserServiceInterface
}

func NewUserController(service services.UserServiceInterface) *UserController {
	return &UserController{
		userService: service,
	}
}

func requestError(reqErr error) rest_errors.RestErr {
	switch{
		case errors.Is(reqErr, user_errors.RequestTimeoutErr):
			return rest_errors.NewRestError("request timeout", http.StatusRequestTimeout, "database error", nil)
		case errors.Is(reqErr, user_errors.NotFoundErr):
			return rest_errors.NewNotFoundError("Ueser not found with given id")
		case errors.Is(reqErr, user_errors.BadRequestErr):
			return rest_errors.NewBadRequestError(reqErr.Error())
		default:
			return rest_errors.NewInternalServerError("internal server error", errors.New("database error"))
	}
}

func getUserId(userIdParam string) (int64, rest_errors.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, rest_errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

func (uc *UserController) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, saveError := uc.userService.CreateUser(ctx, user)
	if saveError != nil {
		restErr := requestError(saveError)
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func (uc *UserController) Get(c *gin.Context) {
	ctx := c.Request.Context()
	if err := oauth.AutenticationRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}
	userId, idErr := getUserId(c.Param("users_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	user, getErr := uc.userService.GetUser(ctx, userId)
	if getErr != nil {
		restErr := requestError(getErr)
		c.JSON(restErr.Status(), restErr)
		return
	}

	if oauth.GetCallerId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func (uc *UserController) Patch(c *gin.Context) {
	ctx := c.Request.Context()
	userId, idErr := getUserId(c.Param("users_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	var user users.PartialUser

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	user.Id = userId

	updatedUser, err := uc.userService.PartialUpdateUser(ctx, user)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, updatedUser.Marshall(c.GetHeader("X-Public") == "true"))
}

func (uc *UserController) Put(c *gin.Context) {
	ctx := c.Request.Context()
	userId, idErr := getUserId(c.Param("users_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	user.Id = userId

	updatedUser, err := uc.userService.UpdateUser(ctx, user)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, updatedUser.Marshall(c.GetHeader("X-Public") == "true"))
}

func (uc *UserController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	userId, idErr := getUserId(c.Param("users_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	if err := uc.userService.DeleteUser(ctx, userId); err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (uc *UserController) Search(c *gin.Context) {
	ctx := c.Request.Context()
	status := c.Query("status")

	users, err := uc.userService.Search(ctx, status)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, users.Marshall(c.GetHeader("X-Public") == "true"))
}

func (uc *UserController) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}
	user, err := uc.userService.LoginUser(ctx, request)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}

	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}
