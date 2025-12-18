package handlers

import (
	"fmt" // Add fmt import
	"io"
	"log/slog"
	"net/http"

	"github.com/GeorgiiMalishev/ideas-platform/config"
	"github.com/GeorgiiMalishev/ideas-platform/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type ImageHandler struct {
	imageUsecase usecase.ImageUsecase
	cfg          *config.Config
	logger       *slog.Logger
}

func NewImageHandler(imageUsecase usecase.ImageUsecase, cfg *config.Config, logger *slog.Logger) *ImageHandler {
	return &ImageHandler{
		imageUsecase: imageUsecase,
		cfg:          cfg,
		logger:       logger,
	}
}

// @Summary Get image
// @Description Get an image by its path stored in MinIO.
// @Tags images
// @Produce octet-stream
// @Param imagePath path string true "Full path of the image (e.g., images/uuid.jpg)"
// @Success 200 {file} byte
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /images/{imagePath} [get]
func (h *ImageHandler) GetImage(c *gin.Context) {
	imagePath := c.Param("imagePath")
	logger := h.logger.With("method", "GetImage", "imagePath", imagePath)
	logger.Debug("attempting to retrieve image")

	// Get the configured bucket name
	configuredBucketName := h.cfg.ImageDB.BucketName

	// Check if the imagePath starts with the configured bucket name and a slash
	// to extract the actual object name within the bucket.
	objectName := imagePath
	if len(imagePath) > len(configuredBucketName) && imagePath[:len(configuredBucketName)] == configuredBucketName && imagePath[len(configuredBucketName)] == '/' {
		objectName = imagePath[len(configuredBucketName)+1:]
	} else if imagePath == configuredBucketName {
		// This case means the imagePath *is* just the bucket name, which is not a valid object name
		logger.Error("invalid image path: image path is just bucket name", slog.String("path", imagePath))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image path"})
		return
	}
	// If it doesn't start with the bucket name, it might be an object in a different bucket
	// or an incorrectly formed path. For this proxy, we assume it's in the configured bucket.


	object, objectInfo, err := h.imageUsecase.GetImage(c.Request.Context(), objectName)
	if err != nil {
		logger.Error("failed to get image from usecase", slog.String("error", err.Error()))
		// Check for specific MinIO errors, e.g., NoSuchKey
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve image"})
		return
	}
	defer object.Close() // Ensure the object stream is closed

	c.Header("Content-Type", objectInfo.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", objectInfo.Size))
	// You might want to add Cache-Control headers here

	if _, err := io.Copy(c.Writer, object); err != nil {
		logger.Error("failed to stream image to client", slog.String("error", err.Error()))
		// Note: Cannot send HTTP status code here as headers have already been written.
		// The client might receive an incomplete file or a connection reset.
		// Logging the error is sufficient.
	}
}
