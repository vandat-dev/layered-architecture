package s3

import (
	"app/global"
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type S3Provider struct {
	client     *minio.Client
	bucketName string
}

func NewS3Provider() *S3Provider {
	return &S3Provider{
		client:     global.MinIO,
		bucketName: global.Config.MinIO.BucketName,
	}
}

// UploadFile uploads a file to S3
func (s *S3Provider) UploadFile(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	info, err := s.client.PutObject(ctx, s.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	global.Logger.Info(fmt.Sprintf("Successfully uploaded %s of size %d", objectName, info.Size))

	// Return the URL (assuming public or presigned)
	// For now, let's return the object name or a constructed URL
	// If using MinIO, the URL might be http://endpoint/bucket/object
	protocol := "http"
	if global.Config.MinIO.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, global.Config.MinIO.Endpoint, s.bucketName, objectName), nil
}

// GetPresignedURL generates a presigned URL for a file
func (s *S3Provider) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presignedURL.String(), nil
}
