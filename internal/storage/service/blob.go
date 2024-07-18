package service

import (
	"context"
	"fmt"
	"medioa/internal/storage/models"
	"medioa/pkg/log"
	"mime/multipart"
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/vukyn/kuery/crypto"
)

// Upload to public Blob Storage with process (handle concurrent chunks)
// https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/storage/azblob/blockblob/examples_test.go
func (s *service) UploadPublicBlob(ctx context.Context, req *models.UploadBlobRequest) (*models.UploadBlobResponse, error) {
	log := log.New("service", "UploadPublicBlob")

	// init new block blob connection
	token := crypto.HashedToken()
	blobName := path.Join("public", token+path.Ext(req.File.Filename))

	// create a request progress object to track the progress of the upload
	totalBytes := req.File.Size
	pr := func(bytesTransferred int64) {
		percentage := float64(bytesTransferred) / float64(totalBytes) * 100
		log.Info("Wrote %d of %d bytes (%.2f%%)\n", bytesTransferred, totalBytes, percentage)
		ws := s.lib.SocketConn.Get(req.SessionId)
		if ws != nil {
			ws.Write([]byte(fmt.Sprintf("%f", percentage)))
		}
	}

	if err := s.uploadBlob(ctx, blobName, req.File, pr); err != nil {
		return nil, err
	}

	return &models.UploadBlobResponse{
		Token:    token,
		FileName: blobName,
		Ext:      path.Ext(req.File.Filename),
		Url:      path.Join(s.cfg.AzBlob.Host, s.cfg.Storage.Container, blobName),
	}, nil
}

// Upload to private Blob Storage with process (handle concurrent chunks)
// https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/storage/azblob/blockblob/examples_test.go
func (s *service) UploadPrivateBlob(ctx context.Context, req *models.UploadBlobRequest) (*models.UploadBlobResponse, error) {
	log := log.New("service", "UploadPrivateBlob")

	if req.SecretId == "" {
		return nil, fmt.Errorf("missing secret id before upload private blob")
	}

	// init new block blob connection
	token := crypto.HashedToken()
	blobName := path.Join("private", req.SecretId, token+path.Ext(req.File.Filename))

	// create a request progress object to track the progress of the upload
	totalBytes := req.File.Size
	pr := func(bytesTransferred int64) {
		percentage := float64(bytesTransferred) / float64(totalBytes) * 100
		log.Info("Wrote %d of %d bytes (%.2f%%)\n", bytesTransferred, totalBytes, percentage)
		ws := s.lib.SocketConn.Get(req.SessionId)
		if ws != nil {
			ws.Write([]byte(fmt.Sprintf("%f", percentage)))
		}
	}

	if err := s.uploadBlob(ctx, blobName, req.File, pr); err != nil {
		return nil, err
	}

	return &models.UploadBlobResponse{
		Token:    token,
		FileName: blobName,
		Ext:      path.Ext(req.File.Filename),
		Url:      path.Join(s.cfg.AzBlob.Host, s.cfg.Storage.Container, blobName),
	}, nil
}

// Download from Blob Storage with SAS (Shared Access Signature)
// https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/storage/azblob/sas/examples_test.go
func (s *service) DownloadSAS(ctx context.Context, req *models.DownloadSASRequest) (*models.DownloadSASResponse, error) {
	log := log.New("service", "DownloadSAS")
	host := s.cfg.AzBlob.Host
	containerName := s.cfg.Storage.Container

	blobURL := fmt.Sprintf("%s/%s/%s", host, containerName, req.FileName)
	blobCli, err := blob.NewClientWithSharedKeyCredential(blobURL, s.lib.Blob.Credential, nil)
	if err != nil {
		log.Error("blockblob.NewClientWithSharedKeyCredential", err)
		return nil, err
	}

	now := time.Now().Add(-10 * time.Second)
	expiry := now.Add(30 * 24 * time.Hour) // 30 days
	permissions := sas.BlobPermissions{Read: true}

	sasURL, err := blobCli.GetSASURL(permissions, expiry, nil)
	if err != nil {
		log.Error("blobCli.GetSASURL", err)
		return nil, err
	}
	fmt.Println(sasURL)

	return &models.DownloadSASResponse{
		Url: sasURL,
	}, nil
}

func (s *service) uploadBlob(ctx context.Context, blobName string, file *multipart.FileHeader, pr func(bytesTransferred int64)) error {
	log := log.New("service", "uploadBlob")

	// open file
	reader, err := file.Open()
	if err != nil {
		log.Error("file.Open", err)
		return err
	}
	defer reader.Close()

	// add progress reporting
	reqProgress := streaming.NewRequestProgress(reader, pr)

	opts := &blockblob.UploadStreamOptions{
		// Metadata: map[string]*string{},
	}
	blobClient := s.lib.Blob.Container.NewBlockBlobClient(blobName)
	if _, err := blobClient.UploadStream(ctx, reqProgress, opts); err != nil {
		log.Error("blobClient.UploadStream", err)
		return err
	}

	return nil
}
