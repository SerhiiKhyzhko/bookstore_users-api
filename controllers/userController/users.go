package userController

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/SerhiiKhyzhko/bookstore-oauth-go/v2/jwtsdk"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/services"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/user_errors"
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
	switch {
	case errors.Is(reqErr, user_errors.RequestTimeoutErr):
		return rest_errors.NewRestError("request timeout", http.StatusRequestTimeout, "database error", nil)
	case errors.Is(reqErr, user_errors.NotFoundErr):
		return rest_errors.NewNotFoundError(fmt.Sprintf("User not found: %s", reqErr.Error()))
	case errors.Is(reqErr, user_errors.BadRequestErr):
		return rest_errors.NewBadRequestError(fmt.Sprintf("Bad request: %s", reqErr.Error()))
	default:
		return rest_errors.NewInternalServerError("internal server error", reqErr)
	}
}

func getFromContext[T any](c *gin.Context, key string) (T, bool) {
	value, exist := c.Get(key)
	if !exist {
		var empty T
		return empty, false
	}

	result, ok := value.(T)
	return result, ok
}

// @Summary     Create new user
// @Description It is a public endpoint. Create new user with provided informtion. If X-Public header is true return 'public user' with
// @Description the most common fields. If X-Public header is false return 'private user' with all fields accept password
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body users.User true "User request"
// @Success     201 {object} users.PrivateUser
// @Failure     400 {object} user_errors.SwaggerRestErr
// @Failure     408 {object} user_errors.SwaggerRestErr
// @Failure     500 {object} user_errors.SwaggerRestErr
// @Router      /users [post]
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

// @Summary     Get user
// @Description It is a private endpoint whitch require jwt token. Return user using id obtained via URL and access token. If X-Public header is true return 'public user' with
// @Description the most common fields. If X-Public header is false return 'private user' with all fields accept password
// @Tags        users
// @Produce     json
// @Param       id path string true "User id"
// @Param       Authorization header string true "Bearer token"
// @Success     200 {object} users.PrivateUser
// @Failure     408 {object} user_errors.SwaggerRestErr
// @Failure     404 {object} user_errors.SwaggerRestErr
// @Failure     500 {object} user_errors.SwaggerRestErr
// @Router      /users [get]
func (uc *UserController) Get(c *gin.Context) {
	ctx := c.Request.Context()
	claims, ok := getFromContext[*jwtsdk.Claims](c, jwtsdk.ClaimsKey)
	if !ok {
		restErr := rest_errors.NewInternalServerError("empty claims", fmt.Errorf("internal server error"))
		c.JSON(restErr.Status(), restErr)
		return
	}
	userId := claims.UserId

	user, getErr := uc.userService.GetUser(ctx, userId)
	if getErr != nil {
		restErr := requestError(getErr)
		c.JSON(restErr.Status(), restErr)
		return
	}
	isPublic, ok := getFromContext[bool](c, jwtsdk.IsPublicKey)
	if !ok {
		isPublic = true // make request public to return short\public user data sturuct
	}
	c.JSON(http.StatusOK, user.Marshall(isPublic))
}

// @Summary     Partially update user fields
// @Description It is a private endpoint whitch require jwt token. Update the user with provided informtion. If X-Public header is true return 'public user' with
// @Description the most common fields. If X-Public header is false return 'private user' with all fields accept password
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body users.PartialUser true "update request"
// @Success     200 {object} users.PrivateUser
// @Failure     400 {object} user_errors.SwaggerRestErr
// @Failure     408 {object} user_errors.SwaggerRestErr
// @Failure     500 {object} user_errors.SwaggerRestErr
// @Router      /users [patch]
func (uc *UserController) Patch(c *gin.Context) {
	ctx := c.Request.Context()
	claims, ok := getFromContext[*jwtsdk.Claims](c, jwtsdk.ClaimsKey)
	if !ok {
		restErr := rest_errors.NewInternalServerError("empty claims", fmt.Errorf("internal server error"))
		c.JSON(restErr.Status(), restErr)
		return
	}
	userId := claims.UserId

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

// @Summary     Replace entire user
// @Description It is a private endpoint whitch require jwt token. Update user with provided informtion. If X-Public header is true return 'public user' with
// @Description the most common fields. If X-Public header is false return 'private user' with all fields accept password
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body users.User true "update request"
// @Success     200 {object} users.PrivateUser
// @Failure     400 {object} user_errors.SwaggerRestErr
// @Failure     408 {object} user_errors.SwaggerRestErr
// @Failure     500 {object} user_errors.SwaggerRestErr
// @Router      /users [put]
func (uc *UserController) Put(c *gin.Context) {
	ctx := c.Request.Context()
	claims, ok := getFromContext[*jwtsdk.Claims](c, jwtsdk.ClaimsKey)
	if !ok {
		restErr := rest_errors.NewInternalServerError("empty claims", fmt.Errorf("internal server error"))
		c.JSON(restErr.Status(), restErr)
		return
	}
	userId := claims.UserId

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

// @Summary     Delete user
// @Description It is a private endpoint whitch require jwt token. Delete an user using id obtained via URL
// @Tags        users
// @Produce     json
// @Param       id path string true "User id"
// @Success     200 {object} map[string]string
// @Failure     404 {object} user_errors.SwaggerRestErr
// @Failure     408 {object} user_errors.SwaggerRestErr
// @Failure     500 {object} user_errors.SwaggerRestErr
// @Router      /users [delete]
func (uc *UserController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	claims, ok := getFromContext[*jwtsdk.Claims](c, jwtsdk.ClaimsKey)
	if !ok {
		restErr := rest_errors.NewInternalServerError("empty claims", fmt.Errorf("internal server error"))
		c.JSON(restErr.Status(), restErr)
		return
	}
	userId := claims.UserId

	if err := uc.userService.DeleteUser(ctx, userId); err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// @Summary     Search users
// @Description It is a private endpoint whitch require jwt token. Returns an array of users If X-Public header is true return 'public user' with
// @Description the most common fields. If X-Public header is false return 'private user' with all fields accept password
// @Tags        users
// @Accept      json
// @Produce     json
// @Param 		status query string false "User status filter"
// @Success     200 {object} []users.User
// @Failure     408 {object} user_errors.SwaggerRestErr
// @Failure     500 {object} user_errors.SwaggerRestErr
// @Router      /internal/users/search [get]
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

// @Summary     Log in user
// @Description It is a public endpoint. Returns user specific user whitch match given email and password. If X-Public
// @Description header is true return 'public user' with the most common fields. If X-Public header
// @Description is false return 'private user' with all fields accept password
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body users.LoginRequest true "Log in request"
// @Success     200 {object} users.PrivateUser
// @Failure     400 {object} user_errors.SwaggerRestErr
// @Failure     408 {object} user_errors.SwaggerRestErr
// @Failure     500 {object} user_errors.SwaggerRestErr
// @Router      /users/login [post]
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
