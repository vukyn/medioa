package handler

import (
	"github.com/gin-gonic/gin"
)

type IHandler interface {
	MapRoutes(group *gin.RouterGroup)
	Upload(ctx *gin.Context)
	UploadWithSecret(ctx *gin.Context)
	Download(ctx *gin.Context)
	DownloadWithSecret(ctx *gin.Context)
	CreateSecret(ctx *gin.Context)
}
