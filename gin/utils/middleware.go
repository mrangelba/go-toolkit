package utils

import (
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/mrangelba/go-toolkit/gin/middlewares/jwtauth"
)

func MiddlewareAuthAndTimeout(auth bool, duration time.Duration, handle gin.HandlerFunc) gin.HandlerFunc {
	if auth && duration > 0 {
		return jwtauth.New(jwtauth.WithHandler(timeout.New(timeout.WithHandler(handle), timeout.WithTimeout(duration))))
	}

	if !auth && duration > 0 {
		return timeout.New(timeout.WithHandler(handle), timeout.WithTimeout(duration))
	}

	if auth && duration == 0 {
		return jwtauth.New(jwtauth.WithHandler(handle))
	}

	return handle
}

func MiddlewareAuth(auth bool, handle gin.HandlerFunc) gin.HandlerFunc {
	if auth {
		return jwtauth.New(jwtauth.WithHandler(handle))
	}

	return handle
}

func MiddlewareTimeout(duration time.Duration, handle gin.HandlerFunc) gin.HandlerFunc {
	if duration > 0 {
		return timeout.New(timeout.WithHandler(handle), timeout.WithTimeout(duration))
	}

	return handle
}
