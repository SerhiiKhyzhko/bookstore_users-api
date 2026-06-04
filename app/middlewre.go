package app

import (
	"errors"
	"net/http"

	"github.com/SerhiiKhyzhko/bookstore-oauth-go/v2/jwtErrors"
	"github.com/SerhiiKhyzhko/bookstore-oauth-go/v2/jwtsdk"
	"github.com/gin-gonic/gin"
)

func authMiddleware(jwtSdk *jwtsdk.JwtManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqType, err := jwtSdk.IsPublic(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}
		c.Set(jwtsdk.IsPublicKey, reqType)

		claims, err := jwtSdk.AuthenticationRequest(c.Request)
		if err != nil {
			if errors.Is(err, jwtErrors.BadRequestErr) {
				c.AbortWithStatusJSON(http.StatusBadRequest, err)
				return
			} else if errors.Is(err, jwtErrors.UnauthorizedErr) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, err)
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, err)
				return
			}
		}
		c.Set(jwtsdk.ClaimsKey, claims)
		c.Next()
	}
}
