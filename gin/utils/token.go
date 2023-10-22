package utils

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/mrangelba/go-toolkit/gin/middlewares/jwtauth"
)

func ExtractTokenClaims(c *gin.Context) (map[string]interface{}, error) {
	token, _ := c.Value(jwtauth.TokenCtxKey).(jwt.Token)

	var err error
	var claims map[string]interface{}

	if token != nil {
		claims, err = token.AsMap(context.Background())
		if err != nil {
			return nil, err
		}
	} else {
		claims = map[string]interface{}{}
	}

	err, _ = c.Value(jwtauth.ErrorCtxKey).(error)

	return claims, err
}

func ExtractTokenSubject(c *gin.Context) (string, error) {
	claims, err := ExtractTokenClaims(c)
	if err != nil {
		return "", err
	}

	return claims["sub"].(string), nil
}
