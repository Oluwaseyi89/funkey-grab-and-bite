package v1

import (
	"net/http"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	handlers.Success(c, gin.H{
		"status":    "ok",
		"service":   "funkey-bite-api",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *HealthHandler) GetHealthLegacy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "funkey-bite-api",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
