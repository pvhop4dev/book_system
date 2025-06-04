package restapi

import (
	"book_system/internal/infrastructure"
	"book_system/internal/model/dto"
	"book_system/internal/service"
	_ "book_system/internal/transport/response"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FileResponse represents the structure of a file upload response
type FileResponse struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type,omitempty"`
	URL      string `json:"url,omitempty"`
}

type uploadController struct {
	uploadService service.IUploadService
}

func NewUploadController(uploadService service.IUploadService) *uploadController {
	return &uploadController{
		uploadService: uploadService,
	}
}

// UploadFile handles single file upload
// @Summary Upload a single file
// @Description Upload a file to the storage
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} FileResponse
// @Router /api/v1/upload [post]
func (u *uploadController) UploadFile(c *gin.Context) {
	// Single file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Validate file size (e.g., 10MB limit)
	const maxUploadSize = 10 << 20 // 10 MB
	if file.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds the limit of 10MB"})
		return
	}

	// Upload the file
	filePath, err := u.uploadService.UploadFile(c.Request.Context(), file, "uploads")
	if err != nil {
		log.Error().Err(err).Str("file", file.Filename).Msg("Failed to upload file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
		return
	}

	// Generate URL for the uploaded file
	fileURL, err := u.uploadService.GetFileURL(c.Request.Context(), filePath)
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to generate file URL")
		// Continue even if URL generation fails, return the file path
	}

	c.JSON(http.StatusOK, dto.FileResponse{
		FileName: filePath,
		URL:      fileURL,
		Size:     file.Size,
	})
}

// UploadMultipleFiles handles multiple file uploads
// @Summary Upload multiple files
// @Description Upload multiple files to the storage
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param files formData []file true "Files to upload"
// @Success 200 {object} response.Response{data=[]FileResponse}
// @Router /api/v1/upload/multiple [post]
func (u *uploadController) UploadMultipleFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files provided"})
		return
	}

	var responses []dto.FileResponse
	ctx := c.Request.Context()

	for _, file := range files {
		// Validate file size (e.g., 10MB limit per file)
		const maxUploadSize = 10 << 20 // 10 MB
		if file.Size > maxUploadSize {
			responses = append(responses, dto.FileResponse{
				FileName: file.Filename,
				URL:      "",
			})
			continue
		}

		// Upload the file
		filePath, err := u.uploadService.UploadFile(ctx, file, "uploads")
		if err != nil {
			log.Error().Err(err).Str("file", file.Filename).Msg("Failed to upload file")
			responses = append(responses, dto.FileResponse{
				FileName: file.Filename,
				URL:      "",
			})
			continue
		}

		// Generate URL for the uploaded file
		fileURL, err := u.uploadService.GetFileURL(ctx, filePath)
		if err != nil {
			log.Error().Err(err).Str("file", filePath).Msg("Failed to generate file URL")
			// Continue even if URL generation fails, return the file path
		}

		responses = append(responses, dto.FileResponse{
			FileName: filePath,
			URL:      fileURL,
			Size:     file.Size,
		})
	}

	c.JSON(http.StatusOK, responses)
}

// GetFile handles file download
// @Summary Get a file
// @Description Get a file by its name
// @Tags files
// @Produce octet-stream
// @Param filename path string true "File name"
// @Success 200 {file} file
// @Router /api/v1/files/{filename} [get]
func (u *uploadController) GetFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	ctx := c.Request.Context()
	file, err := infrastructure.GetFile(ctx, filename)
	if err != nil {
		log.Error().Err(err).Str("filename", filename).Msg("Failed to get file")
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	defer file.Close()

	// Get file info for content type
	objInfo, err := file.Stat()
	if err != nil {
		log.Error().Err(err).Str("filename", filename).Msg("Failed to get file info")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file info"})
		return
	}

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", objInfo.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", objInfo.Size))

	// Stream the file to the response
	c.Stream(func(w io.Writer) bool {
		_, err := io.Copy(w, file)
		if err != nil {
			log.Error().Err(err).Str("filename", filename).Msg("Failed to stream file")
			return false
		}
		return false
	})
}

// DeleteFile handles file deletion
// @Summary Delete a file
// @Description Delete a file by its name
// @Tags files
// @Produce json
// @Param filename path string true "File name"
// @Success 200 {object} map[string]string
// @Router /api/v1/files/{filename} [delete]
func (u *uploadController) DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	ctx := c.Request.Context()
	err := infrastructure.DeleteFile(ctx, filename)
	if err != nil {
		log.Error().Err(err).Str("filename", filename).Msg("Failed to delete file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted successfully"})
}

// GetFileURL generates a pre-signed URL for a file
// @Summary Get a pre-signed URL for a file
// @Description Generate a pre-signed URL for a file
// @Tags files
// @Produce json
// @Param filename path string true "File name"
// @Param expiry query int false "Expiry time in hours (default: 24)"
// @Success 200 {object} map[string]string
// @Router /api/v1/files/{filename}/url [get]
func (u *uploadController) GetFileURL(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	expiry := 24 * time.Hour
	if expiryStr := c.Query("expiry"); expiryStr != "" {
		if hours, err := time.ParseDuration(expiryStr + "h"); err == nil {
			expiry = hours
		}
	}

	ctx := c.Request.Context()
	url, err := infrastructure.GeneratePresignedURL(ctx, filename, expiry)
	if err != nil {
		log.Error().Err(err).Str("filename", filename).Msg("Failed to generate file URL")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate file URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":    url.String(),
		"expiry": expiry.String(),
	})
}
