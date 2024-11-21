package handler

import (
	"fmt"
	"medioa/config"
	"medioa/constants"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

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
	group.POST(constants.STORAGE_ENDPOINT_UPLOAD_STAGE, h.UploadChunk)
	group.POST(constants.STORAGE_ENDPOINT_UPLOAD_COMMIT, h.CommitChunk)
	group.POST(constants.STORAGE_ENDPOINT_UPLOAD_WITH_SECRET, h.UploadWithSecret)
	group.POST(constants.STORAGE_ENDPOINT_UPLOAD_STAGE_WITH_SECRET, h.UploadChunkWithSecret)
	group.POST(constants.STORAGE_ENDPOINT_UPLOAD_COMMIT_WITH_SECRET, h.CommitChunkWithSecret)
	group.GET(constants.STORAGE_ENDPOINT_DOWNLOAD, h.Download)
	group.GET(constants.STORAGE_ENDPOINT_REQUEST_DOWNLOAD, h.RequestDownload)
	group.POST(constants.STORAGE_ENDPOINT_CREATE_SECRET, h.CreateSecret)
	group.PUT(constants.STORAGE_ENDPOINT_RETRIEVE_SECRET, h.RetrieveSecret)
	group.PUT(constants.STORAGE_ENDPOINT_RESET_PIN_CODE, h.ResetPinCode)
}

// Upload godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Upload media
//	@Description	Upload media file (images, videos, etc.), must provide file or url
//	@Tags			Storage
//	@Accept			mpfd
//	@Produce		multipart/form-data
//	@Param			id			query		string	false	"session id"
//	@Param			url			formData	string	false	"file url"
//	@Param			file		formData	file	false	"binary file"
//	@Param			file_name	formData	string	false	"file name"
//	@Success		201			{object}	models.UploadResponse
//	@Router			/storage/upload [post]
func (h Handler) Upload(ctx *gin.Context) {
	maxSize := h.cfg.Upload.MaxSizeMB
	contentType := ctx.GetHeader("Content-Type")
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize<<20)

	id := ctx.Query("id")
	fileName := ctx.PostForm("file_name")
	url := ctx.PostForm("url")
	var file *multipart.FileHeader
	if strings.Contains(contentType, "multipart/form-data") {
		var err error
		file, err = ctx.FormFile("file")
		if err != nil {
			if err.Error() != "http: no such file" {
				if err.Error() == "multipart: NextPart: http: request body too large" {
					xhttp.BadRequest(ctx, fmt.Errorf("file size too large (max: %dMB)", maxSize))
				} else {
					xhttp.BadRequest(ctx, err)
				}
				return
			}
		}
	}
	userId := int64(1)

	req := &models.UploadRequest{
		SessionId: id,
		URL:       url,
		File:      file,
		FileName:  fileName,
	}

	if err := req.Validate(); err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	res, err := h.usecase.Upload(ctx, userId, req)
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Created(ctx, res)
}

// UploadChunk godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Upload media by chunk
//	@Description	Upload media file (images, videos, etc.)
//	@Tags			Storage
//	@Accept			mpfd
//	@Produce		json
//	@Param			id				query		string	false	"session id"
//	@Param			chunk			formData	file	true	"binary chunk"
//	@Param			chunk_index		formData	int64	true	"chunk index"
//	@Param			total_chunks	formData	int64	true	"total chunk"
//	@Param			file_id			formData	string	false	"file id"
//	@Param			file_name		formData	string	false	"file name"
//	@Success		201				{object}	models.UploadChunkResponse
//	@Router			/storage/upload/stage [post]
func (h Handler) UploadChunk(ctx *gin.Context) {
	maxSize := h.cfg.Upload.MaxSizeMB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize<<20)

	id := ctx.Query("id")
	fileId := ctx.PostForm("file_id")
	fileName := ctx.PostForm("file_name")
	chunk, err := ctx.FormFile("chunk")
	if err != nil {
		if err.Error() == "multipart: NextPart: http: request body too large" {
			xhttp.BadRequest(ctx, fmt.Errorf("chunk size too large (max: %dMB)", maxSize))
		} else {
			xhttp.BadRequest(ctx, err)
		}
		return
	}

	chunkIndexStr := ctx.PostForm("chunk_index")
	chunkIndex, err := strconv.ParseInt(chunkIndexStr, 10, 64)
	if err != nil {
		xhttp.BadRequest(ctx, fmt.Errorf("invalid chunk index"))
		return
	}

	totalChunkStr := ctx.PostForm("total_chunks")
	totalChunks, err := strconv.ParseInt(totalChunkStr, 10, 64)
	if err != nil {
		xhttp.BadRequest(ctx, fmt.Errorf("invalid total chunk"))
		return
	}

	userId := int64(1)
	res, err := h.usecase.UploadChunk(ctx, userId, &models.UploadChunkRequest{
		SessionId:   id,
		FileId:      fileId,
		FileName:    fileName,
		Chunk:       chunk,
		ChunkIndex:  chunkIndex,
		TotalChunks: totalChunks,
	})
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Created(ctx, res)
}

// CommitChunk godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Commit upload media chunk
//	@Description	Commit all chunks to complete upload media file
//	@Tags			Storage
//	@Accept			json
//	@Produce		json
//	@Param			id		query		string						false	"session id"
//	@Param			body	body		models.CommitChunkRequest	true	"commit chunk request"
//	@Success		200		{object}	models.CommitChunkResponse
//	@Router			/storage/upload/commit [post]
func (h Handler) CommitChunk(ctx *gin.Context) {
	id := ctx.Query("id")
	req := &models.CommitChunkRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}
	req.SessionId = id

	userId := int64(1)
	res, err := h.usecase.CommitChunk(ctx, userId, req)
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Created(ctx, res)
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
//	@Param			file		formData	file	true	"binary file"
//	@Param			file_name	formData	string	false	"file name"
//	@Success		201			{object}	models.UploadResponse
//	@Router			/storage/secret/upload [post]
func (h Handler) UploadWithSecret(ctx *gin.Context) {
	id := ctx.Query("id")
	secret := ctx.Query("secret")
	fileName := ctx.PostForm("file_name")
	file, err := ctx.FormFile("file")
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

// UploadChunkWithSecret godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Upload media by chunk with secret
//	@Description	Upload media file (images, videos, etc.)
//	@Tags			Storage
//	@Accept			mpfd
//	@Produce		json
//	@Param			id				query		string	false	"session id"
//	@Param			secret			query		string	true	"secret"
//	@Param			chunk			formData	file	true	"binary chunk"
//	@Param			chunk_index		formData	int64	true	"chunk index"
//	@Param			total_chunks	formData	int64	true	"total chunk"
//	@Param			file_id			formData	string	false	"file id"
//	@Param			file_name		formData	string	false	"file name"
//	@Success		201				{object}	models.UploadChunkResponse
//	@Router			/storage/secret/upload/stage [post]
func (h Handler) UploadChunkWithSecret(ctx *gin.Context) {
	maxSize := h.cfg.Upload.MaxSizeMB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize<<20)

	id := ctx.Query("id")
	secret := ctx.Query("secret")
	fileId := ctx.PostForm("file_id")
	fileName := ctx.PostForm("file_name")
	chunk, err := ctx.FormFile("chunk")
	if err != nil {
		if err.Error() == "multipart: NextPart: http: request body too large" {
			xhttp.BadRequest(ctx, fmt.Errorf("chunk size too large (max: %dMB)", maxSize))
		} else {
			xhttp.BadRequest(ctx, err)
		}
		return
	}

	chunkIndexStr := ctx.PostForm("chunk_index")
	chunkIndex, err := strconv.ParseInt(chunkIndexStr, 10, 64)
	if err != nil {
		xhttp.BadRequest(ctx, fmt.Errorf("invalid chunk index"))
		return
	}

	totalChunkStr := ctx.PostForm("total_chunks")
	totalChunks, err := strconv.ParseInt(totalChunkStr, 10, 64)
	if err != nil {
		xhttp.BadRequest(ctx, fmt.Errorf("invalid total chunk"))
		return
	}

	userId := int64(1)
	res, err := h.usecase.UploadChunkWithSecret(ctx, userId, &models.UploadChunkWithSecretRequest{
		SessionId:   id,
		Secret:      secret,
		FileId:      fileId,
		FileName:    fileName,
		Chunk:       chunk,
		ChunkIndex:  chunkIndex,
		TotalChunks: totalChunks,
	})
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Created(ctx, res)
}

// CommitChunkWithSecret godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Commit upload media chunk with secret
//	@Description	Commit all chunks to complete upload media file
//	@Tags			Storage
//	@Accept			json
//	@Produce		json
//	@Param			id		query		string						false	"session id"
//	@Param			secret	query		string						true	"secret"
//	@Param			body	body		models.CommitChunkRequest	true	"commit chunk request"
//	@Success		200		{object}	models.CommitChunkResponse
//	@Router			/storage/secret/upload/commit [post]
func (h Handler) CommitChunkWithSecret(ctx *gin.Context) {
	id := ctx.Query("id")
	secret := ctx.Query("secret")
	req := &models.CommitChunkRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}
	req.SessionId = id
	req.Secret = secret

	userId := int64(1)
	res, err := h.usecase.CommitChunkWithSecret(ctx, userId, req)
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	xhttp.Created(ctx, res)
}

// RequestDownload godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Request download private media
//	@Description	Get download url for private media with download password
//	@Tags			Storage
//	@Accept			json
//	@Produce		json
//	@Param			file_id	path		string	true	"file id"
//	@Param			token	query		string	true	"token"
//	@Param			secret	query		string	true	"secret"
//	@Success		200		{object}	models.RequestDownloadResponse
//	@Router			/storage/download/request/{file_id} [get]
func (h Handler) RequestDownload(ctx *gin.Context) {
	userId := int64(1)
	fileId := ctx.Param("file_id")
	token := ctx.Query("token")
	secret := ctx.Query("secret")
	res, err := h.usecase.RequestDownload(ctx, userId, &models.RequestDownloadRequest{
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

// Download godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Download media (public/private)
//	@Description	Download media file (images, videos, etc.)
//	@Tags			Storage
//	@Accept			json
//	@Produce		json
//	@Param			file_id		path		string	true	"file id"
//	@Param			token		query		string	true	"token"
//	@Param			secret		query		string	false	"secret"
//	@Param			password	query		string	false	"password"
//	@Param			silent		query		bool	false	"silent response"
//	@Success		200			{object}	models.DownloadResponse
//	@Router			/storage/download/{file_id} [get]
func (h Handler) Download(ctx *gin.Context) {
	userId := int64(1)
	fileId := ctx.Param("file_id")
	token := ctx.Query("token")
	secret := ctx.Query("secret")
	password := ctx.Query("password")
	silent := ctx.Query("silent")
	res, err := h.usecase.Download(ctx, userId, &models.DownloadRequest{
		FileId:           fileId,
		Token:            token,
		Secret:           secret,
		DownloadPassword: password,
	})
	if err != nil {
		xhttp.BadRequest(ctx, err)
		return
	}

	if silent != "true" {
		xhttp.Redirect(ctx, res.Url)
	} else {
		xhttp.Ok(ctx, res)
	}
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
