package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/requestctx"
)

type fakeAuditRecorder struct {
	entries []models.AuditEntry
	actors  []requestctx.Actor
}

func (f *fakeAuditRecorder) Record(ctx context.Context, exec repository.Execer, entry models.AuditEntry) error {
	f.entries = append(f.entries, entry)
	f.actors = append(f.actors, requestctx.ActorFrom(ctx))
	return nil
}

func TestAuditWritesRecordsSuccessfulAdminChanges(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := &fakeAuditRecorder{}
	router := gin.New()
	admin := router.Group("/api/v1/admin", func(c *gin.Context) {
		c.Set("user_id", 3)
		c.Set("is_admin", true)
	}, AuditWrites(recorder, "POST /api/v1/admin/orders/:id/payment-review"))
	admin.PATCH("/orders/:id/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })
	admin.GET("/orders/:id", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })
	admin.PUT("/menu/items/:id", func(c *gin.Context) { c.JSON(http.StatusBadRequest, gin.H{}) })
	admin.POST("/orders/:id/payment-review", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })

	for _, r := range []struct{ method, path string }{
		{http.MethodPatch, "/api/v1/admin/orders/17/status"},
		{http.MethodGet, "/api/v1/admin/orders/17"},
		{http.MethodPut, "/api/v1/admin/menu/items/4"},
		{http.MethodPost, "/api/v1/admin/orders/17/payment-review"},
	} {
		router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(r.method, r.path, nil))
	}

	if len(recorder.entries) != 1 {
		t.Fatalf("only the successful status change should be recorded, got %+v", recorder.entries)
	}
	entry := recorder.entries[0]
	if entry.Action != "admin.patch /api/v1/admin/orders/:id/status" || entry.EntityType != "orders" || entry.EntityID != "17" {
		t.Fatalf("unexpected entry %+v", entry)
	}
	if actor := recorder.actors[0]; actor.Type != requestctx.ActorAdmin || actor.ID != "3" {
		t.Fatalf("expected admin 3 as the actor, got %+v", actor)
	}
}
