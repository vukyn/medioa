package service

import (
	"context"
	"medioa/internal/storage/models"
	"medioa/pkg/log"
	"path"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/vukyn/kuery/crypto"
)

// Upload to Blob Storage with process (handle concurrent with large file)
// https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/storage/azblob/blockblob/examples_test.go
func (s *service) Upload(ctx context.Context, req *models.UploadRequest) (*models.UploadResponse, error) {
	log := log.New("service", "Upload")

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
	blobClient := s.lib.BlobContainer.NewBlockBlobClient(blobName)

	// create a request progress object to track the progress of the upload
	totalBytes := req.File.Size
	reqProgress := streaming.NewRequestProgress(reader, func(bytesTransferred int64) {
		percent := float64(bytesTransferred) / float64(totalBytes) * 100
		log.Info("Wrote %d of %d bytes (%.2f%%)\n", bytesTransferred, totalBytes, percent)
		socketConn := s.lib.SocketConn.Get(req.SessionId)
		if socketConn != nil {
			socketConn.Emit("upload_progress", map[string]any{
				"token": token,
				"percent": percent,
			})
		}
	})

	opts := &blockblob.UploadStreamOptions{
		// Metadata: map[string]*string{},
	}
	if _, err := blobClient.UploadStream(ctx, reqProgress, opts); err != nil {
		log.Error("blobClient.UploadStream", err)
		return nil, err
	}

	return &models.UploadResponse{
		Token:    token,
		FileName: blobName,
		Ext:      path.Ext(req.File.Filename),
		Url:      path.Join(s.cfg.AzBlob.Host, s.cfg.Storage.Container, blobName),
	}, nil
}
