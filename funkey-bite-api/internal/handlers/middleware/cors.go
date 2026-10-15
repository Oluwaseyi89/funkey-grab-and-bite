package middleware

import (
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// defaultAllowedOrigins are always allowed; CORS_ALLOWED_ORIGINS adds to them.
var defaultAllowedOrigins = []string{
	"http://localhost:3000",
	"http://localhost:5173",
	"http://localhost:3001",
	"https://funkeygrabandbite.com",
}

// CORSMiddleware echoes back allowed origins. Extra origins come from the
// comma-separated CORS_ALLOWED_ORIGINS env var; an entry may contain a single
// "*" wildcard (e.g. https://my-app-*.vercel.app for Vercel preview
// deployments) which matches one DNS label segment, never a "." or "/".
func CORSMiddleware() gin.HandlerFunc {
	exactOrigins := map[string]bool{}
	var wildcardOrigins []string
	for _, origin := range append(defaultAllowedOrigins, strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")...) {
		origin = strings.TrimRight(strings.TrimSpace(origin), "/")
		switch {
		case origin == "":
		case strings.Count(origin, "*") == 1:
			wildcardOrigins = append(wildcardOrigins, origin)
		default:
			exactOrigins[origin] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin != "" {
			if exactOrigins[origin] || isLocalOrigin(origin) || matchesWildcardOrigin(origin, wildcardOrigins) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
			c.Writer.Header().Add("Vary", "Origin")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods",
			"POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// isLocalOrigin allows any port on localhost/127.0.0.1 for local frontend dev.
func isLocalOrigin(origin string) bool {
	parsed, err := url.Parse(origin)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return false
	}
	host := parsed.Hostname()
	return host == "localhost" || host == "127.0.0.1"
}

func matchesWildcardOrigin(origin string, patterns []string) bool {
	for _, pattern := range patterns {
		prefix, suffix, _ := strings.Cut(pattern, "*")
		if len(origin) <= len(prefix)+len(suffix) ||
			!strings.HasPrefix(origin, prefix) || !strings.HasSuffix(origin, suffix) {
			continue
		}
		if middle := origin[len(prefix) : len(origin)-len(suffix)]; !strings.ContainsAny(middle, "./:") {
			return true
		}
	}
	return false
}
