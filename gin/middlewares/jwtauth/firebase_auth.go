package jwtauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FirebaseAuthVerify(cfg *config) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := tokenFromHeader(c.Request)
		if bearerToken == "" {
			c.AbortWithError(http.StatusUnauthorized, ErrNoTokenFound)
			return
		}

		token, err := cfg.firebaseAuth.VerifyIDToken(c.Request.Context(), bearerToken)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set(TokenCtxKey, token)
		c.Set(ErrorCtxKey, err)

		if token != nil {
			c.Set(UserCtxKey, token.Subject)
		}

		if cfg.handler != nil {
			cfg.handler(c)
		}

		c.Next()
	}
}
