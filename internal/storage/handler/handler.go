package handler

import (
	"medioa/config"
	"medioa/constants"

	"medioa/internal/storage/models"
	"medioa/internal/storage/usecase"
	commonModel "medioa/models"
	"medioa/pkg/http"
	"medioa/pkg/log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg     *config.Config
	lib     *commonModel.Lib
	usecase usecase.IUsecase
}

func InitHandler(cfg *config.Config, lib *commonModel.Lib, usecase usecase.IUsecase) IHandler {
	return Handler{
		cfg:     cfg,
		lib:     lib,
		usecase: usecase,
	}
}

func (h Handler) MapRoutes(group *gin.RouterGroup) {
	group.POST(constants.STORAGE_ENDPOINT_UPLOAD, h.UploadMedia)
}

// UploadMedia godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Upload media
//	@Description	Upload media file (images, videos, etc.)
//	@Tags			Storage
//	@Accept			mpfd
//	@Produce		json
//	@Param			chunk	formData	file	true	"binary file"
//	@Success		201		{object}	models.UploadResponse
//	@Router			/storage/upload [post]
func (h Handler) UploadMedia(ctx *gin.Context) {
	log := log.New("handler", "UploadMedia")

	file, err := ctx.FormFile("chunk")
	if err != nil {
		log.Error("Failed to get file from request", err)
		http.BadRequest(ctx, err)
		return
	}

	userId := int64(1)
	res, err := h.usecase.Upload(ctx, userId, &models.UploadRequest{
		File: file,
	})
	if err != nil {
		http.Internal(ctx, err)
		return
	}

	http.Created(ctx, res)
}
