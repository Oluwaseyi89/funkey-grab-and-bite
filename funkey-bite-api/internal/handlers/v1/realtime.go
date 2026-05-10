package v1

import (
	"net/http"
	"strings"

	"funkey-grab-and-bite/funkey-bite-api/internal/realtime"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type RealtimeHandler struct {
	adminRepo repository.IAdminRepository
	upgrader  websocket.Upgrader
}

func NewRealtimeHandler(adminRepo repository.IAdminRepository) *RealtimeHandler {
	return &RealtimeHandler{
		adminRepo: adminRepo,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
	}
}

func (h *RealtimeHandler) ConnectAdmin(c *gin.Context) {
	tokenString := strings.TrimSpace(c.Query("token"))
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token is required"})
		return
	}

	claims, err := utils.ValidateToken(tokenString)
	if err != nil || !claims.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	admin, err := h.adminRepo.GetAdminUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify admin access"})
		return
	}
	if admin == nil || !admin.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	realtime.GlobalHub.RegisterConnection(conn)
}
