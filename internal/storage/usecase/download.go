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

	if params.Secret != "" {
		// get secret info
		secret, err := u.verifySecretToken(ctx, params.Secret)
		if err != nil {
			return nil, err
		}

		// secret invalid
		if file.SecretId != secret.UUID {
			return nil, fmt.Errorf("permission denied")
		}
	} else {
		// secret required
		if file.SecretId != "" {
			// no request download found
			if file.DownloadPassword == "" {
				return nil, fmt.Errorf("permission denied")
			}
		}
	}

	// password invalid
	if params.DownloadPassword != file.DownloadPassword {
		return nil, fmt.Errorf("permission denied")
	}

	//
	if file.CreatedBy != userId {
		return nil, fmt.Errorf("permission denied")
	}

	// end validation

	var sas *azBlobModel.DownloadSASResponse
	downloadFileName := map[bool]string{true: file.Token + file.Ext, false: file.Token}[file.Ext != ""]

	if file.SecretId == "" {
		// public download
		downloadFileName = path.Join("public", downloadFileName)
		sas, err = u.azBlobSv.DownloadSAS(ctx, &azBlobModel.DownloadSASRequest{
			FileName: downloadFileName,
		})
		if err != nil {
			log.Error("usecase.azBlobSv.DownloadSAS", err)
			return nil, err
		}
	} else {
		// private download
		downloadFileName = path.Join("private", file.SecretId, downloadFileName)
		sas, err = u.azBlobSv.DownloadSAS(ctx, &azBlobModel.DownloadSASRequest{
			FileName: downloadFileName,
		})
		if err != nil {
			log.Error("usecase.azBlobSv.DownloadSAS", err)
			return nil, err
		}
	}

	return &storageModel.DownloadResponse{
		Url: sas.Url,
	}, nil
}

func (u *usecase) RequestDownload(ctx context.Context, userId int64, params *storageModel.RequestDownloadRequest) (*storageModel.RequestDownloadResponse, error) {
	log := log.New("usecase", "RequestDownload")

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

	// end validation

	downloadPassword := file.DownloadPassword
	if downloadPassword == "" {
		downloadPassword = generateDownloadPassword()
		if _, err := u.storageSv.Update(ctx, userId, &storageModel.SaveRequest{
			UUID:             file.UUID,
			DownloadPassword: downloadPassword,
		}); err != nil {
			log.Error("usecase.storageSv.Update", err)
			return nil, err
		}
	}

	return &storageModel.RequestDownloadResponse{
		Url:      getDownloadUrl(u.cfg.App.Host, file.UUID, file.Token),
		Password: downloadPassword,
		FileName: file.FileName,
	}, nil
}
