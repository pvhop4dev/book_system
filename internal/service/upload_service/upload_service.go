package uploadservice

import (
	"book_system/internal/infrastructure"
	"context"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
)

type uploadService struct {
}

func NewUploadService() *uploadService {
	return &uploadService{}
}

func (s *uploadService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, customPath ...string) (string, error) {
	return infrastructure.UploadFile(ctx, fileHeader, customPath...)
}

func (s *uploadService) DeleteFile(ctx context.Context, objectName string) error {
	return infrastructure.DeleteFile(ctx, objectName)
}

func (s *uploadService) GetFileURL(ctx context.Context, objectName string) (string, error) {
	return infrastructure.GetFileURL(ctx, objectName)
}

func (s *uploadService) GetFile(ctx context.Context, objectName string) (*multipart.FileHeader, error) {
	minioObj, err := infrastructure.GetFile(ctx, objectName)
	if err != nil {
		return nil, err
	}

	// Create a temporary file to store the MinIO object's content
	tempFile, err := os.CreateTemp("", "minio-obj-"+objectName)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Copy MinIO object content to the temporary file
	_, err = io.Copy(tempFile, minioObj)
	if err != nil {
		return nil, err
	}

	// Create a multipart.FileHeader from the temporary file
	// Get object info to get size and content type
	objInfo, err := minioObj.Stat()
	if err != nil {
		return nil, err
	}

	fileHeader := &multipart.FileHeader{
		Filename: filepath.Base(objectName),
		Size:     objInfo.Size,
	}
	fileHeader.Header = textproto.MIMEHeader{}
	fileHeader.Header.Add("Content-Type", objInfo.ContentType)

	return fileHeader, nil
}
