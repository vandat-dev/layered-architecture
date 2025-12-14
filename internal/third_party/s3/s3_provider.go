package s3

import (
	"app/global"
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"time"

	"go.uber.org/zap"

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
	if contentType == "" || contentType == "application/octet-stream" {
		// Try to detect from filename extension
		ext := filepath.Ext(objectName)
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		default:
			contentType = "application/octet-stream"
		}
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

// ListObjects lists objects in a specific prefix
func (s *S3Provider) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	var objects []string
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}

	for object := range s.client.ListObjects(ctx, s.bucketName, opts) {
		if object.Err != nil {
			return nil, object.Err
		}
		// Construct URL
		protocol := "http"
		if global.Config.MinIO.UseSSL {
			protocol = "https"
		}
		url := fmt.Sprintf("%s://%s/%s/%s", protocol, global.Config.MinIO.Endpoint, s.bucketName, object.Key)
		objects = append(objects, url)
	}
	return objects, nil
}

// RemoveFolder removes all objects with a specific prefix
func (s *S3Provider) RemoveFolder(ctx context.Context, prefix string) error {
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		opts := minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		}
		for object := range s.client.ListObjects(ctx, s.bucketName, opts) {
			if object.Err != nil {
				global.Logger.Error("Error listing objects for deletion", zap.Error(object.Err))
				return
			}
			objectsCh <- object
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	for err := range s.client.RemoveObjects(ctx, s.bucketName, objectsCh, opts) {
		if err.Err != nil {
			return fmt.Errorf("failed to remove object %s: %w", err.ObjectName, err.Err)
		}
	}
	return nil
}

func (s *S3Provider) UploadBytes(ctx context.Context, data []byte, objectName string) (string, error) {
	reader := bytes.NewReader(data)
	size := int64(len(data))
	contentType := "image/webp" // Default for this use case, or pass as arg

	info, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload bytes: %w", err)
	}

	global.Logger.Info(fmt.Sprintf("Successfully uploaded bytes %s of size %d", objectName, info.Size))

	protocol := "http"
	if global.Config.MinIO.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, global.Config.MinIO.Endpoint, s.bucketName, objectName), nil
}
