package handler

import (
	"github.com/gin-gonic/gin"
)

type IHandler interface {
	MapRoutes(group *gin.RouterGroup)
}
