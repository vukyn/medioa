package usecase

import (
	"context"
	"fmt"
	"medioa/constants"
	"medioa/pkg/xtype"
	"net/url"
	"path"
	"strings"

	"github.com/vukyn/kuery/log"

	secretModel "medioa/internal/secret/models"
	storageModel "medioa/internal/storage/models"

	"github.com/vukyn/kuery/cryp"
	"github.com/zRedShift/mimemagic"
)

func generateDownloadPassword() string {
	return cryp.HashUUID()[0:8]
}

func sniffMimeType(file xtype.File) (string, error) {
	log := log.New("usecase", "sniffMimeType")

	reader, err := file.Open()
	if err != nil {
		log.Error("file.Open", err)
		return "", err
	}
	defer reader.Close()

	mimeType, err := mimemagic.MatchReader(reader, file.Filename)
	if err != nil {
		log.Error("mimemagic.MatchReader", err)
		return "", err
	}

	return mimeType.MediaType(), nil
}

func getUploadedFileName1(file xtype.File) string {
	var fileName string
	ext := path.Ext(file.Filename)
	fileName = strings.ReplaceAll(file.Filename, ext, "")
	if fileName == "" {
		fileName = file.Filename
	}
	return fileName
}

func getUploadedFileName2(URL string) string {
	var fileName string

	url, err := url.Parse(URL)
	if err != nil {
		return ""
	}

	fileName = path.Base(url.Path)
	ext := path.Ext(fileName)
	fileName = strings.ReplaceAll(fileName, ext, "")
	if fileName == "" {
		fileName = path.Base(url.Path)
	}
	return fileName
}

func getDownloadUrl(host, fileId, token string) string {
	filePath := strings.ReplaceAll(constants.SHARE_ENDPOINT_DOWNLOAD, ":file_id", fileId)
	downloadUrl := fmt.Sprintf("%s/share%s?token=%s", host, filePath, token)
	return downloadUrl
}

func (u *usecase) verifySecretToken(ctx context.Context, secretToken string) (*secretModel.Response, error) {
	log := log.New("usecase", "verifySecretToken")

	if secretToken == "" {
		return nil, fmt.Errorf("secret token is required")
	}

	secret, err := u.secretSv.GetOne(ctx, &secretModel.RequestParams{
		AccessToken: secretToken,
	})
	if err != nil {
		log.Error("usecase.secretSv.GetOne", err)
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("secret token is invalid")
	}

	return secret, nil
}

func (u *usecase) verifyFileInfo(ctx context.Context, fileId, token string) (*storageModel.Response, error) {
	log := log.New("usecase", "verifyFileInfo")

	if fileId == "" {
		return nil, fmt.Errorf("file id is required")
	}

	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	file, err := u.storageSv.GetOne(ctx, &storageModel.RequestParams{
		UUID:  fileId,
		Token: token,
	})
	if err != nil {
		log.Error("usecase.storageSv.GetOne", err)
		return nil, err
	}
	if file == nil {
		return nil, fmt.Errorf("file not found")
	}

	return file, nil
}

func (u *usecase) getFileById(ctx context.Context, fileId string) (*storageModel.Response, error) {
	log := log.New("usecase", "getFileById")

	if fileId == "" {
		return nil, fmt.Errorf("file id is required")
	}

	file, err := u.storageSv.GetOne(ctx, &storageModel.RequestParams{
		UUID: fileId,
	})
	if err != nil {
		log.Error("usecase.storageSv.GetOne", err)
		return nil, err
	}
	if file == nil {
		return nil, fmt.Errorf("file not found")
	}

	return file, nil
}
