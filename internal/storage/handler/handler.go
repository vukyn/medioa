package handler

import (
	"fmt"
	"medioa/config"
	"medioa/constants"
	"net/http"

	"medioa/internal/storage/models"
	"medioa/internal/storage/usecase"
	commonModel "medioa/models"
	"medioa/pkg/xhttp"

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
	group.POST(constants.STORAGE_ENDPOINT_UPLOAD_WITH_SECRET, h.UploadWithSecret)
	group.GET(constants.STORAGE_ENDPOINT_DOWNLOAD_WITH_SECRET, h.DownloadWithSecret)
	group.POST(constants.STORAGE_ENDPOINT_CREATE_SECRET, h.CreateSecret)
	group.PUT(constants.STORAGE_ENDPOINT_RETRIEVE_SECRET, h.RetrieveSecret)
	group.PUT(constants.STORAGE_ENDPOINT_RESET_PIN_CODE, h.ResetPinCode)
}

// Upload godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Upload media
//	@Description	Upload media file (images, videos, etc.)
//	@Tags			Storage
//	@Accept			mpfd
//	@Produce		json
//	@Param			id			query		string	false	"session id"
//	@Param			chunk		formData	file	true	"binary file"
//	@Param			file_name	formData	string	false	"file name"
//	@Success		201			{object}	models.UploadResponse
//	@Router			/storage/upload [post]
func (h Handler) Upload(ctx *gin.Context) {
	maxSize := h.cfg.Upload.MaxSizeMB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize<<20)

	id := ctx.Query("id")
	fileName := ctx.PostForm("file_name")
	file, err := ctx.FormFile("chunk")
	if err != nil {
		if err.Error() == "multipart: NextPart: http: request body too large" {
			xhttp.BadRequest(ctx, fmt.Errorf("file size too large (max: %dMB)", maxSize))
		} else {
			xhttp.BadRequest(ctx, err)
		}
		return
	}
	userId := int64(1)
	res, err := h.usecase.Upload(ctx, userId, &models.UploadRequest{
		SessionId: id,
		File:      file,
		FileName:  fileName,
	})
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Created(ctx, res)
}

// Download godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Download media
//	@Description	Download media file
//	@Tags			Storage
//	@Accept			json
//	@Produce		json
//	@Param			file_id	path		string	true	"file id"
//	@Param			token	query		string	true	"token"
//	@Success		200		{object}	models.DownloadResponse
//	@Router			/storage/download/{file_name} [get]
func (h Handler) Download(ctx *gin.Context) {
	userId := int64(1)
	fileId := ctx.Param("file_id")
	token := ctx.Query("token")
	res, err := h.usecase.Download(ctx, userId, &models.DownloadRequest{
		FileId: fileId,
		Token:  token,
	})
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Ok(ctx, res)
}

// UploadWithSecret godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Upload media with secret
//	@Description	Upload media file (images, videos, etc.)
//	@Tags			Storage
//	@Accept			mpfd
//	@Produce		json
//	@Param			id			query		string	false	"session id"
//	@Param			secret		query		string	true	"secret"
//	@Param			chunk		formData	file	true	"binary file"
//	@Param			filename	formData	string	false	"file name"
//	@Success		201			{object}	models.UploadResponse
//	@Router			/storage/secret/upload [post]
func (h Handler) UploadWithSecret(ctx *gin.Context) {
	id := ctx.Query("id")
	secret := ctx.Query("secret")
	fileName := ctx.PostForm("filename")
	file, err := ctx.FormFile("chunk")
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	userId := int64(1)
	res, err := h.usecase.UploadWithSecret(ctx, userId, &models.UploadWithSecretRequest{
		SessionId: id,
		Secret:    secret,
		File:      file,
		FileName:  fileName,
	})
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Created(ctx, res)
}

// DownloadWithSecret godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Download media with secret
//	@Description	Download media file
//	@Tags			Storage
//	@Accept			json
//	@Produce		json
//	@Param			file_id	path		string	true	"file id"
//	@Param			token	query		string	true	"token"
//	@Param			secret	query		string	true	"secret"
//	@Success		200		{object}	models.DownloadResponse
//	@Router			/storage/secret/download/{file_name} [get]
func (h Handler) DownloadWithSecret(ctx *gin.Context) {
	userId := int64(1)
	fileId := ctx.Param("file_id")
	token := ctx.Query("token")
	secret := ctx.Query("secret")
	res, err := h.usecase.DownloadWithSecret(ctx, userId, &models.DownloadWithSecretRequest{
		FileId: fileId,
		Token:  token,
		Secret: secret,
	})
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Ok(ctx, res)
}

// CreateSecret godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Create new secret
//	@Description	Create new secret for upload media
//	@Tags			Storage
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.CreateSecretRequest	true	"create secret request"
//	@Success		201		{object}	models.CreateSecretResponse
//	@Router			/storage/secret [post]
func (h Handler) CreateSecret(ctx *gin.Context) {
	userId := int64(1)
	req := &models.CreateSecretRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}
	res, err := h.usecase.CreateSecret(ctx, userId, req)
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Created(ctx, res)
}

// RetrieveSecret godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Retrieve secret
//	@Description	Retrieve secret with new access token
//	@Tags			Storage
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.RetrieveSecretRequest	true	"retrieve secrect request"
//	@Success		200		{object}	models.RetrieveSecretResponse
//	@Router			/storage/secret/retrieve [put]
func (h Handler) RetrieveSecret(ctx *gin.Context) {
	userId := int64(1)
	req := &models.RetrieveSecretRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}
	res, err := h.usecase.RetrieveSecret(ctx, userId, req)
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Ok(ctx, res)
}

// ResetPinCode godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Reset pin code
//	@Description	Reset pin code for secret
//	@Tags			Storage
//	@Accept			json
//	@Produce		json
//	@Param			body	body	models.ResetPinCodeRequest	true	"reset pin request"
//	@Success		200
//	@Router			/storage/secret/pin [put]
func (h Handler) ResetPinCode(ctx *gin.Context) {
	userId := int64(1)
	req := &models.ResetPinCodeRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}
	res, err := h.usecase.ResetPinCode(ctx, userId, req)
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Ok(ctx, res)
}
