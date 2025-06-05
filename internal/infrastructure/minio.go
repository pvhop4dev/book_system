package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"sync"
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
	initMinioOnce sync.Once
)

func InitMinio() error {
	var err error
	initMinioOnce.Do(func() {
		err = initMinio()
		if err != nil {
			log.Error().Err(err).Msg("MinIO client initialization failed, will retry later")
		}
	})
	return err
}
func GetMinioClient() *minio.Client {
	return minioClient
}

func CloseMinio() {
}

// InitMinio initializes the MinIO client with configuration from viper
func initMinio() error {
	cfg := config.MustGet()

	// Initialize minio client object
	endpoint := cfg.Minio.Host + ":" + cfg.Minio.Port
	accessKeyID := cfg.Minio.AccessKey
	secretAccessKey := cfg.Minio.SecretKey
	useSSL := cfg.Minio.Secure

	// Thêm retry mechanism
	maxRetries := 10
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			log.Info().Msgf("Attempt %d to connect to MinIO", i+1)
			time.Sleep(retryDelay)
		}

		client, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create MinIO client")
			continue
		}

		// Kiểm tra bucket
		bucketExists, err := client.BucketExists(context.Background(), cfg.Minio.DefaultBucket)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check bucket existence")
			continue
		}

		if !bucketExists {
			if err := client.MakeBucket(context.Background(), cfg.Minio.DefaultBucket, minio.MakeBucketOptions{}); err != nil {
				log.Error().Err(err).Msg("Failed to create bucket")
				continue
			}
			log.Info().Msgf("Bucket '%s' created successfully", cfg.Minio.DefaultBucket)
		}

		// Set bucket policy if needed
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		policy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::` + cfg.Minio.DefaultBucket + `/*"]
			}]
		}`

		if err := client.SetBucketPolicy(ctx, cfg.Minio.DefaultBucket, policy); err != nil {
			log.Error().Err(err).Msg("Failed to set bucket policy")
			continue
		}

		// Lưu client và cấu hình
		minioClient = client
		defaultBucket = cfg.Minio.DefaultBucket
		returnURL = cfg.Minio.ReturnURL
		log.Info().Msgf("MinIO client initialized successfully with bucket: %s", cfg.Minio.DefaultBucket)
		return nil
	}

	return fmt.Errorf("failed to initialize MinIO client after %d attempts", maxRetries)
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
		return errors.New("MinIO client not initialized")
	}

	bucket := GetDefaultBucket()
	_, err := minioClient.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil // File doesn't exist, consider this a success
		}
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	err = minioClient.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
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

// GetDefaultBucket returns the default bucket name
func GetDefaultBucket() string {
	return defaultBucket
}

// GetFileURL generates a URL for a single file
func GetFileURL(ctx context.Context, objectName string) (string, error) {
	url, err := GeneratePresignedURL(ctx, objectName, 24*time.Hour)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
