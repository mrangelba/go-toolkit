package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrangelba/go-toolkit/errors"
	"github.com/mrangelba/go-toolkit/problem_details"
)

func CatchErrorProblemDetails(c *gin.Context, err error) {
	if err == errors.ErrInvalidRequest {
		c.AbortWithStatusJSON(http.StatusBadRequest, problem_details.NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if err == errors.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusNotFound, problem_details.NewHTTPError(http.StatusNotFound, err))
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, problem_details.NewHTTPError(http.StatusInternalServerError, err))
		return
	}
}
