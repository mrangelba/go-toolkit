package utils

import (
	"time"

	"firebase.google.com/go/v4/auth"
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

func MiddlewareAuth(handle gin.HandlerFunc) gin.HandlerFunc {
	return jwtauth.New(jwtauth.WithHandler(handle))
}

func MiddlewareFirebaseAuth(firebaseAuth *auth.Client, handle gin.HandlerFunc) gin.HandlerFunc {
	return jwtauth.New(jwtauth.WithFirebaseAuth(firebaseAuth), jwtauth.WithHandler(handle))
}

func MiddlewareTimeout(duration time.Duration, handle gin.HandlerFunc) gin.HandlerFunc {
	return timeout.New(timeout.WithHandler(handle), timeout.WithTimeout(duration))
}
