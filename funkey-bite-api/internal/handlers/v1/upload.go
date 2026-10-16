package v1

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"

	"funkey-grab-and-bite/funkey-bite-api/internal/storage"

	"github.com/gin-gonic/gin"
)

// MaxMenuImageBytes caps a menu image upload (matches the admin form's limit).
const MaxMenuImageBytes = 2 << 20

// allowedImageTypes maps sniffed content types to the extension stored in the key.
var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type UploadHandler struct {
	storage storage.Storage
}

// NewUploadHandler accepts a nil storage; uploads then answer 503 so the rest
// of the API keeps working while storage is not configured.
func NewUploadHandler(s storage.Storage) *UploadHandler {
	return &UploadHandler{storage: s}
}

// UploadMenuImage stores a menu image and returns its key and public URL
// @Summary Upload menu image
// @Description Upload a JPEG, PNG or WebP image (max 2 MB) for a menu item. Save the returned key as the item's imageUrl.
// @Tags admin
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Image file"
// @Success 201 {object} map[string]string
// @Router /admin/uploads/menu-image [post]
func (h *UploadHandler) UploadMenuImage(c *gin.Context) {
	if h.storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Image storage is not configured"})
		return
	}

	// Bound the whole request body; the multipart envelope adds a little
	// overhead on top of the file itself.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxMenuImageBytes+64<<10)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Image must be 2 MB or smaller"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "An image file is required in the \"file\" field"})
		return
	}
	if fileHeader.Size > MaxMenuImageBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Image must be 2 MB or smaller"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read the uploaded file"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, MaxMenuImageBytes+1))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read the uploaded file"})
		return
	}
	if len(data) > MaxMenuImageBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Image must be 2 MB or smaller"})
		return
	}

	// Trust the bytes, not the client-supplied filename or Content-Type.
	contentType := http.DetectContentType(data)
	extension, ok := allowedImageTypes[contentType]
	if !ok {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Only JPEG, PNG and WebP images are allowed"})
		return
	}

	key, err := storage.NewMenuImageKey(extension)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store image"})
		return
	}

	if err := h.storage.Put(c.Request.Context(), key, contentType, bytes.NewReader(data), int64(len(data))); err != nil {
		log.Printf("menu image upload failed: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to store image"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"key": key,
		"url": storage.ResolveURL(h.storage, key),
	})
}
