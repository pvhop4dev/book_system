package infrastructure

import (
	"context"
	"errors"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"time"

	"book_system/internal/config"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

var (
	minioClient   *minio.Client
	defaultBucket string
	returnURL     string
)

// InitMinio initializes the MinIO client with configuration from viper
func InitMinio() error {
	cfg := config.MustGet()

	// Initialize minio client object
	endpoint := cfg.Minio.Host + ":" + cfg.Minio.Port
	accessKeyID := cfg.Minio.AccessKey
	secretAccessKey := cfg.Minio.SecretKey
	useSSL := cfg.Minio.Secure

	// Initialize minio client object
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	// Set global variables
	minioClient = client
	defaultBucket = cfg.Minio.DefaultBucket
	returnURL = cfg.Minio.ReturnURL

	// Create default bucket if not exists
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, defaultBucket)
	if err != nil {
		return err
	}

	if !exists {
		err = client.MakeBucket(ctx, defaultBucket, minio.MakeBucketOptions{
			Region:        cfg.Minio.Location,
			ObjectLocking: false,
		})
		if err != nil {
			return err
		}

		// Set bucket policy if needed
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + defaultBucket + `/*"]}]}`
		err = client.SetBucketPolicy(ctx, defaultBucket, policy)
		if err != nil {
			return err
		}
	}

	log.Info().Str("bucket", defaultBucket).Msg("MinIO client initialized successfully")
	return nil
}

// UploadFile uploads a file to MinIO
func UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, customPath ...string) (string, error) {
	if minioClient == nil {
		return "", errors.New("minio client not initialized")
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Generate a unique filename if needed
	fileExt := filepath.Ext(fileHeader.Filename)
	objectName := uuid.New().String() + fileExt
	if len(customPath) > 0 && customPath[0] != "" {
		objectName = customPath[0] + "/" + objectName
	}

	// Upload the file to MinIO
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = minioClient.PutObject(
		ctx,
		defaultBucket,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType:  contentType,
			UserMetadata: map[string]string{"original-filename": fileHeader.Filename},
		},
	)
	if err != nil {
		return "", err
	}

	// Generate URL for the uploaded file
	fileURL := returnURL + "/" + defaultBucket + "/" + objectName
	return fileURL, nil
}

// GetFile retrieves a file from MinIO
func GetFile(ctx context.Context, objectName string) (*minio.Object, error) {
	if minioClient == nil {
		return nil, errors.New("minio client not initialized")
	}

	obj, err := minioClient.GetObject(
		ctx,
		defaultBucket,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// DeleteFile deletes a file from MinIO
func DeleteFile(ctx context.Context, objectName string) error {
	if minioClient == nil {
		return errors.New("minio client not initialized")
	}

	// Check if the object exists first
	_, err := minioClient.StatObject(ctx, defaultBucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return err
	}

	// Delete the object
	err = minioClient.RemoveObject(ctx, defaultBucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

// GeneratePresignedURL generates a pre-signed URL for a file
func GeneratePresignedURL(ctx context.Context, objectName string, expiry time.Duration) (*url.URL, error) {
	if minioClient == nil {
		return nil, errors.New("minio client not initialized")
	}

	// Set request parameters for content disposition and type.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+objectName)

	// Generates a presigned url which expires in 7 days.
	presignedURL, err := minioClient.PresignedGetObject(
		ctx,
		defaultBucket,
		objectName,
		expiry,
		reqParams,
	)
	if err != nil {
		return nil, err
	}

	return presignedURL, nil
}

// GetMinioClient returns the MinIO client instance
func GetMinioClient() *minio.Client {
	return minioClient
}

// GetDefaultBucket returns the default bucket name
func GetDefaultBucket() string {
	return defaultBucket
}

// Close closes the MinIO client connection
func Close() error {
	// MinIO client doesn't have a close method
	// This is just a placeholder for consistency with other infrastructure components
	return nil
}
