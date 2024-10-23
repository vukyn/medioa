package handler

import (
	"medioa/config"
	"medioa/constants"

	initStorage "medioa/internal/storage/init"
	storageModel "medioa/internal/storage/models"
	storageUC "medioa/internal/storage/usecase"
	commonModel "medioa/models"
	"medioa/pkg/xhttp"

	ratelimiter "github.com/vukyn/kuery/middleware/gin/rate_limiter"

	"github.com/dustin/go-humanize"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg       *config.Config
	lib       *commonModel.Lib
	storageUC storageUC.IUsecase
}

func InitHandler(cfg *config.Config, lib *commonModel.Lib, storage *initStorage.Init) IHandler {
	return Handler{
		cfg:       cfg,
		lib:       lib,
		storageUC: storage.Usecase,
	}
}

func (h Handler) MapRoutes(group *gin.RouterGroup) {
	group.GET(constants.SHARE_ENDPOINT_DOWNLOAD, ratelimiter.LimitPerSecond(constants.RATE_LIMIT_DOWNLOAD_PER_SECOND), h.Download)
}

// Download godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Download media (public/private)
//	@Description	Download media file (images, videos, etc.)
//	@Tags			Share
//	@Accept			json
//	@Produce		json
//	@Param			file_id	path		string	true	"file id"
//	@Param			token	query		string	true	"token"
//	@Success		200		{object}	storageModel.DownloadResponse
//	@Router			/share/download/{file_id} [get]
func (h Handler) Download(ctx *gin.Context) {
	userId := int64(1)
	fileId := ctx.Param("file_id")
	token := ctx.Query("token")
	res, err := h.storageUC.GetFileInfo(ctx, userId, &storageModel.GetFileInfoRequest{
		FileId: fileId,
	})
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.HTML(ctx, "tmpl.share.html", gin.H{
		"token":         token,
		"file_id":       fileId,
		"file_name":     res.FileName,
		"file_size":     res.FileSize,
		"file_size_str": humanize.Bytes(uint64(res.FileSize)),
		"has_secret":    res.HasSecret,
	})
}
