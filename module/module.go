package module

import "github.com/gin-gonic/gin"

type Module interface {
	Register(r *gin.Engine)
}
