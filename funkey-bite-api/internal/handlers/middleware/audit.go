package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

// AuditRecorder appends to the audit log.
type AuditRecorder interface {
	Record(ctx context.Context, exec repository.Execer, entry models.AuditEntry) error
}

// AuditWrites records every successful write on the routes it guards: who
// made it, which route and record, and from where. Request bodies aren't
// stored, since some carry passwords. Routes listed in detailed already write
// a fuller entry (with the state before and after) in their own transaction,
// so they aren't recorded twice.
func AuditWrites(recorder AuditRecorder, detailed ...string) gin.HandlerFunc {
	skip := make(map[string]bool, len(detailed))
	for _, route := range detailed {
		skip[route] = true
	}

	return func(c *gin.Context) {
		c.Next()

		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			return
		}
		status := c.Writer.Status()
		route := routeOf(c)
		if status >= 400 || skip[c.Request.Method+" "+route] {
			return
		}

		entry := models.AuditEntry{
			Action:     "admin." + strings.ToLower(c.Request.Method) + " " + route,
			EntityType: entityTypeOf(route),
			EntityID:   c.Param("id"),
			Metadata:   map[string]any{"status": status},
		}
		if err := recorder.Record(RequestContext(c), nil, entry); err != nil {
			log.Printf("[audit] %s: %v", entry.Action, err)
		}
	}
}

// entityTypeOf names what a route changes, from its first segment after
// /admin/, e.g. "orders" for /api/v1/admin/orders/:id/status.
func entityTypeOf(route string) string {
	rest := route
	if i := strings.Index(route, "/admin/"); i >= 0 {
		rest = route[i+len("/admin/"):]
	}
	if segment, _, _ := strings.Cut(strings.Trim(rest, "/"), "/"); segment != "" {
		return segment
	}
	return "admin"
}
