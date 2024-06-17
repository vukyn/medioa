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
	group.POST(constants.STORAGE_ENDPOINT_UPLOAD, h.Upload)
	group.GET(constants.STORAGE_ENDPOINT_DOWNLOAD, h.Download)
}

// Upload godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Upload media
//	@Description	Upload media file (images, videos, etc.)
//	@Tags			Storage
//	@Accept			mpfd
//	@Produce		json
//	@Param			id		query		string	false	"session id"
//	@Param			chunk	formData	file	true	"binary file"
//	@Success		201		{object}	models.UploadResponse
//	@Router			/storage/upload [post]
func (h Handler) Upload(ctx *gin.Context) {
	log := log.New("handler", "Upload")

	id := ctx.Query("id")
	file, err := ctx.FormFile("chunk")
	if err != nil {
		log.Error("Failed to get file from request", err)
		http.BadRequest(ctx, err)
		return
	}

	userId := int64(1)
	res, err := h.usecase.Upload(ctx, userId, &models.UploadRequest{
		SessionId: id,
		File:      file,
	})
	if err != nil {
		http.Internal(ctx, err)
		return
	}

	http.Created(ctx, res)
}

// Download godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Download media
//	@Description	Download media file
//	@Tags			Storage
//	@Accept			aplication/json
//	@Produce		json
//	@Param			file_name	path		string	true	"file name"
//	@Param			token		query		string	true	"token"
//	@Success		200			{object}	models.DownloadResponse
//	@Router			/storage/download/{file_name} [get]
func (h Handler) Download(ctx *gin.Context) {
	// log := log.New("handler", "Download")

	userId := int64(1)
	fileName := ctx.Param("file_name")
	token := ctx.Query("token")
	res, err := h.usecase.Download(ctx, userId, &models.DownloadRequest{
		FileName: fileName,
		Token:    token,
	})
	if err != nil {
		http.Internal(ctx, err)
		return
	}

	http.Ok(ctx, res)
}

// UploadWithSecret godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Upload media with secret
//	@Description	Upload media file (images, videos, etc.)
//	@Tags			Storage
//	@Accept			mpfd
//	@Produce		json
//	@Param			id		query		string	false	"session id"
//	@Param			secret		query		string	true	"secret"
//	@Param			chunk	formData	file	true	"binary file"
//	@Success		201		{object}	models.UploadResponse
//	@Router			/storage/upload [post]
func (h Handler) UploadWithSecret(ctx *gin.Context) {
	log := log.New("handler", "Upload")

	id := ctx.Query("id")
	file, err := ctx.FormFile("chunk")
	if err != nil {
		log.Error("Failed to get file from request", err)
		http.BadRequest(ctx, err)
		return
	}

	userId := int64(1)
	res, err := h.usecase.Upload(ctx, userId, &models.UploadRequest{
		SessionId: id,
		File:      file,
	})
	if err != nil {
		http.Internal(ctx, err)
		return
	}

	http.Created(ctx, res)
}

// DownloadWithSecret godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Download media with secret
//	@Description	Download media file
//	@Tags			Storage
//	@Accept			aplication/json
//	@Produce		json
//	@Param			file_name	path		string	true	"file name"
//	@Param			token		query		string	true	"token"
//	@Param			secret		query		string	true	"secret"
//	@Success		200			{object}	models.DownloadResponse
//	@Router			/storage/download/{file_name} [get]
func (h Handler) DownloadWithSecret(ctx *gin.Context) {
	// log := log.New("handler", "Download")

	userId := int64(1)
	fileName := ctx.Param("file_name")
	token := ctx.Query("token")
	res, err := h.usecase.Download(ctx, userId, &models.DownloadRequest{
		FileName: fileName,
		Token:    token,
	})
	if err != nil {
		http.Internal(ctx, err)
		return
	}

	http.Ok(ctx, res)
}

// CreateSecret godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Create new secret
//	@Description	Create new secret for upload media
//	@Tags			Storage
//	@Accept			mpfd
//	@Produce		json
//	@Param			id		query		string	false	"session id"
//	@Param			secret		query		string	true	"secret"
//	@Param			chunk	formData	file	true	"binary file"
//	@Success		201		{object}	models.UploadResponse
//	@Router			/storage/upload [post]
func (h Handler) CreateSecret(ctx *gin.Context) {
	log := log.New("handler", "Upload")

	id := ctx.Query("id")
	file, err := ctx.FormFile("chunk")
	if err != nil {
		log.Error("Failed to get file from request", err)
		http.BadRequest(ctx, err)
		return
	}

	userId := int64(1)
	res, err := h.usecase.Upload(ctx, userId, &models.UploadRequest{
		SessionId: id,
		File:      file,
	})
	if err != nil {
		http.Internal(ctx, err)
		return
	}

	http.Created(ctx, res)
}