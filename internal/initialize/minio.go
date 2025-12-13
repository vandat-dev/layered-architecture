package initialize

import (
	"app/global"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// InitMinIO initializes the MinIO client
func InitMinIO() {
	global.Logger.Info("Initializing MinIO...")

	endpoint := global.Config.MinIO.Endpoint
	accessKeyID := global.Config.MinIO.AccessKeyID
	secretAccessKey := global.Config.MinIO.SecretAccessKey
	useSSL := global.Config.MinIO.UseSSL
	bucketName := global.Config.MinIO.BucketName

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Failed to initialize MinIO client: %v", err))
		return
	}

	global.MinIO = minioClient

	// Check if bucket exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Failed to check if bucket exists: %v", err))
		return
	}

	if !exists {
		global.Logger.Info(fmt.Sprintf("Bucket %s does not exist, creating it...", bucketName))
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			global.Logger.Error(fmt.Sprintf("Failed to create bucket: %v", err))
			return
		}
		global.Logger.Info(fmt.Sprintf("Bucket %s created successfully", bucketName))
	} else {
		global.Logger.Info(fmt.Sprintf("Bucket %s already exists", bucketName))
	}
}
