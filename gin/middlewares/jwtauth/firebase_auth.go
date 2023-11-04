package jwtauth

import (
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func FirebaseAuthVerify(client *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := tokenFromHeader(c.Request)
		if bearerToken == "" {
			c.AbortWithError(http.StatusUnauthorized, ErrNoTokenFound)
			return
		}

		token, err := client.VerifyIDToken(c.Request.Context(), bearerToken)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set(TokenCtxKey, token)
		c.Set(ErrorCtxKey, err)
		c.Set(UserCtxKey, token.Subject)

		c.Next()
	}
}
