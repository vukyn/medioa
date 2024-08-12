package usecase

import (
	"context"
	"fmt"
	azBlobModel "medioa/internal/azblob/models"
	storageModel "medioa/internal/storage/models"
	"medioa/pkg/log"
	"path"
)

func (u *usecase) Download(ctx context.Context, userId int64, params *storageModel.DownloadRequest) (*storageModel.DownloadResponse, error) {
	log := log.New("usecase", "Download")

	// validation

	// get file info
	file, err := u.verifyFileInfo(ctx, params.FileId, params.Token)
	if err != nil {
		return nil, err
	}

	// check permission
	if file.SecretId != "" {
		return nil, fmt.Errorf("permission denied")
	}
	if file.CreatedBy != userId {
		return nil, fmt.Errorf("permission denied")
	}

	// end validation

	downloadFileName := map[bool]string{true: file.Token + file.Ext, false: file.Token}[file.Ext != ""]
	downloadFileName = path.Join("public", downloadFileName)
	sas, err := u.azBlobSv.DownloadSAS(ctx, &azBlobModel.DownloadSASRequest{
		FileName: downloadFileName,
	})
	if err != nil {
		log.Error("usecase.azBlobSv.DownloadSAS", err)
		return nil, err
	}

	return &storageModel.DownloadResponse{
		Url: sas.Url,
	}, nil
}

func (u *usecase) DownloadWithSecret(ctx context.Context, userId int64, params *storageModel.DownloadWithSecretRequest) (*storageModel.DownloadResponse, error) {
	log := log.New("usecase", "DownloadWithSecret")

	// validation

	// get file info
	file, err := u.verifyFileInfo(ctx, params.FileId, params.Token)
	if err != nil {
		return nil, err
	}

	// get secret info
	secret, err := u.verifySecretToken(ctx, params.Secret)
	if err != nil {
		return nil, err
	}

	// check permission
	if file.SecretId != secret.UUID {
		return nil, fmt.Errorf("permission denied")
	}
	if file.CreatedBy != userId {
		return nil, fmt.Errorf("permission denied")
	}

	// end validation

	downloadFileName := map[bool]string{true: file.Token + file.Ext, false: file.Token}[file.Ext != ""]
	downloadFileName = path.Join("private", file.SecretId, downloadFileName)
	sas, err := u.azBlobSv.DownloadSAS(ctx, &azBlobModel.DownloadSASRequest{
		FileName: downloadFileName,
	})
	if err != nil {
		log.Error("usecase.azBlobSv.DownloadSAS", err)
		return nil, err
	}

	return &storageModel.DownloadResponse{
		Url: sas.Url,
	}, nil
}
