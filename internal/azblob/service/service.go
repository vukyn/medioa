package service

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"medioa/config"
	"medioa/internal/azblob/models"
	commonModel "medioa/models"
	"medioa/pkg/log"
	"medioa/pkg/xtype"
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/vukyn/kuery/crypto"
)

type service struct {
	cfg *config.Config
	lib *commonModel.Lib
}

func InitService(cfg *config.Config, lib *commonModel.Lib) IService {
	return &service{
		cfg: cfg,
		lib: lib,
	}
}

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

	// upload blob
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

// Upload to public Blob Storage by chunk
// https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/storage/azblob/blockblob/examples_test.go
func (s *service) UploadPublicChunk(ctx context.Context, req *models.UploadChunkRequest) (*models.UploadChunkResponse, error) {
	log := log.New("service", "UploadPublicChunk")

	// init new block blob connection
	token := req.Token
	if token == "" {
		token = crypto.HashedToken()
	}
	blobName := path.Join("public", token+path.Ext(req.FileName))

	// open file
	reader, err := req.Chunk.Open()
	if err != nil {
		log.Error("Chunk.Open", err)
		return nil, err
	}
	defer reader.Close()

	// stage block
	blockId := blockIdBase64(req.ChunkIndex)
	opts := &blockblob.StageBlockOptions{}
	blobClient := s.lib.Blob.Container.NewBlockBlobClient(blobName)
	if _, err := blobClient.StageBlock(ctx, blockId, reader, opts); err != nil {
		log.Error("blobClient.StageBlock", err)
		return nil, err
	}

	// progress reporting
	ws := s.lib.SocketConn.Get(req.SessionId)
	if ws != nil {
		percentage := float64(req.ChunkIndex+1) / float64(req.TotalChunks) * 100
		if percentage > 100 {
			percentage = 100
		}
		fmt.Println(percentage)
		ws.Write([]byte(fmt.Sprintf("%f", percentage)))
	}

	return &models.UploadChunkResponse{
		Token:    token,
		BlockId:  blockId,
		FileName: blobName,
		Ext:      path.Ext(req.FileName),
		Url:      path.Join(s.cfg.AzBlob.Host, s.cfg.Storage.Container, blobName),
	}, nil
}

// Commit all public chunks to Blob Storage
func (s *service) CommitPublicChunk(ctx context.Context, req *models.CommitChunkRequest) (*models.CommitChunkRsponse, error) {
	log := log.New("service", "CommitPublicChunk")

	if req.Token == "" {
		return nil, fmt.Errorf("missing token before commit chunk")
	}
	if req.FileName == "" {
		return nil, fmt.Errorf("missing file name before commit chunk")
	}
	if len(req.BlockIds) == 0 {
		return nil, fmt.Errorf("missing block ids before commit chunk")
	}

	// init new block blob connection
	opts := &blockblob.CommitBlockListOptions{}
	blobName := path.Join("public", req.Token+path.Ext(req.FileName))
	blobClient := s.lib.Blob.Container.NewBlockBlobClient(blobName)
	if _, err := blobClient.CommitBlockList(ctx, req.BlockIds, opts); err != nil {
		log.Error("blobClient.CommitBlockList", err)
		return nil, err
	}

	var fileSize int64
	var totalBlock int64
	getBlock, err := blobClient.GetBlockList(ctx, blockblob.BlockListTypeCommitted, nil)
	if err != nil {
		log.Error("blobClient.GetBlockList", err)
		return nil, err
	}
	for _, block := range getBlock.BlockList.CommittedBlocks {
		fileSize += *block.Size
		totalBlock++
	}

	return &models.CommitChunkRsponse{
		FileSize:   fileSize,
		TotalBlock: totalBlock,
	}, nil
}

// Upload to private Blob Storage by chunk
// https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/storage/azblob/blockblob/examples_test.go
func (s *service) UploadPrivateChunk(ctx context.Context, req *models.UploadChunkRequest) (*models.UploadChunkResponse, error) {
	log := log.New("service", "UploadPrivateChunk")

	if req.SecretId == "" {
		return nil, fmt.Errorf("missing secret id before upload private chunk")
	}

	// init new block blob connection
	token := req.Token
	if token == "" {
		token = crypto.HashedToken()
	}
	blobName := path.Join("private", req.SecretId, token+path.Ext(req.FileName))

	// open file
	reader, err := req.Chunk.Open()
	if err != nil {
		log.Error("Chunk.Open", err)
		return nil, err
	}
	defer reader.Close()

	// stage block
	blockId := blockIdBase64(req.ChunkIndex)
	opts := &blockblob.StageBlockOptions{}
	blobClient := s.lib.Blob.Container.NewBlockBlobClient(blobName)
	if _, err := blobClient.StageBlock(ctx, blockId, reader, opts); err != nil {
		log.Error("blobClient.StageBlock", err)
		return nil, err
	}

	// progress reporting
	ws := s.lib.SocketConn.Get(req.SessionId)
	if ws != nil {
		percentage := float64(req.ChunkIndex+1) / float64(req.TotalChunks) * 100
		if percentage > 100 {
			percentage = 100
		}
		fmt.Println(percentage)
		ws.Write([]byte(fmt.Sprintf("%f", percentage)))
	}

	return &models.UploadChunkResponse{
		Token:    token,
		BlockId:  blockId,
		FileName: blobName,
		Ext:      path.Ext(req.FileName),
		Url:      path.Join(s.cfg.AzBlob.Host, s.cfg.Storage.Container, blobName),
	}, nil
}

// Commit all private chunks to Blob Storage
func (s *service) CommitPrivateChunk(ctx context.Context, req *models.CommitChunkRequest) (*models.CommitChunkRsponse, error) {
	log := log.New("service", "CommitPrivateChunk")

	if req.Token == "" {
		return nil, fmt.Errorf("missing token before commit chunk")
	}
	if req.FileName == "" {
		return nil, fmt.Errorf("missing file name before commit chunk")
	}
	if req.SecretId == "" {
		return nil, fmt.Errorf("missing secret id before commit chunk")
	}
	if len(req.BlockIds) == 0 {
		return nil, fmt.Errorf("missing block ids before commit chunk")
	}

	// init new block blob connection
	opts := &blockblob.CommitBlockListOptions{}
	blobName := path.Join("private", req.SecretId, req.Token+path.Ext(req.FileName))
	blobClient := s.lib.Blob.Container.NewBlockBlobClient(blobName)
	if _, err := blobClient.CommitBlockList(ctx, req.BlockIds, opts); err != nil {
		log.Error("blobClient.CommitBlockList", err)
		return nil, err
	}

	var fileSize int64
	var totalBlock int64
	getBlock, err := blobClient.GetBlockList(ctx, blockblob.BlockListTypeCommitted, nil)
	if err != nil {
		log.Error("blobClient.GetBlockList", err)
		return nil, err
	}
	for _, block := range getBlock.BlockList.CommittedBlocks {
		fileSize += *block.Size
		totalBlock++
	}

	return &models.CommitChunkRsponse{
		FileSize:   fileSize,
		TotalBlock: totalBlock,
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

func (s *service) uploadBlob(ctx context.Context, blobName string, file xtype.File, pr func(bytesTransferred int64)) error {
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

func blockIdBase64(idx int64) string {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, idx)
	return base64.StdEncoding.EncodeToString(buf)
}
