package service

import (
	"context"
	"fmt"
	"medioa/internal/storage/models"
	"medioa/pkg/log"
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/vukyn/kuery/crypto"
)

// Upload to Blob Storage with process (handle concurrent with large file)
// https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/storage/azblob/blockblob/examples_test.go
func (s *service) UploadBlob(ctx context.Context, req *models.UploadBlobRequest) (*models.UploadBlobResponse, error) {
	log := log.New("service", "UploadBlob")

	// open file
	reader, err := req.File.Open()
	if err != nil {
		log.Error("req.File.Open", err)
		return nil, err
	}
	defer reader.Close()

	// init new block blob connection
	token := crypto.HashedToken()
	blobName := token + path.Ext(req.File.Filename)
	blobClient := s.lib.Blob.Container.NewBlockBlobClient(blobName)

	// create a request progress object to track the progress of the upload
	totalBytes := req.File.Size
	reqProgress := streaming.NewRequestProgress(reader, func(bytesTransferred int64) {
		percentage := float64(bytesTransferred) / float64(totalBytes) * 100
		log.Info("Wrote %d of %d bytes (%.2f%%)\n", bytesTransferred, totalBytes, percentage)
		ws := s.lib.SocketConn.Get(req.SessionId)
		if ws != nil {
			ws.Write([]byte(fmt.Sprintf("%f", percentage)))
		}
	})

	opts := &blockblob.UploadStreamOptions{
		// Metadata: map[string]*string{},
	}
	if _, err := blobClient.UploadStream(ctx, reqProgress, opts); err != nil {
		log.Error("blobClient.UploadStream", err)
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
	expiry := now.Add(15 * time.Minute)
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
