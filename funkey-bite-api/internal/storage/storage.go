// Package storage uploads static files (menu images) to S3-compatible object
// storage. Supabase Storage and AWS S3 both speak the S3 API, so a single
// implementation serves either; STORAGE_PROVIDER only selects which settings
// are required. The database stores object keys (e.g. "menu-items/ab12.webp")
// and PublicURL turns them into full URLs, so switching providers is a config
// change plus copying the bucket, never a data migration.
package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

const (
	ProviderSupabase = "supabase"
	ProviderS3       = "s3"

	// MenuImagePrefix is the key "folder" every menu image is stored under.
	MenuImagePrefix = "menu-items/"
	// TeamPhotoPrefix is the key "folder" for team member photos.
	TeamPhotoPrefix = "team/"
)

// managedPrefixes are the folders this API uploads into.
var managedPrefixes = []string{MenuImagePrefix, TeamPhotoPrefix}

// Storage is the minimal surface the API needs from an object store.
type Storage interface {
	Put(ctx context.Context, key, contentType string, body io.Reader, size int64) error
	Delete(ctx context.Context, key string) error
	// PublicBaseURL is the URL that object keys are appended to, without a
	// trailing slash.
	PublicBaseURL() string
}

// Config holds the STORAGE_* settings.
type Config struct {
	Provider        string
	Endpoint        string
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	PublicBaseURL   string
}

// ConfigFromEnv reads the STORAGE_* environment variables.
func ConfigFromEnv() Config {
	get := func(key string) string { return strings.TrimSpace(os.Getenv(key)) }
	return Config{
		Provider:        strings.ToLower(get("STORAGE_PROVIDER")),
		Endpoint:        strings.TrimRight(get("STORAGE_ENDPOINT"), "/"),
		Region:          get("STORAGE_REGION"),
		Bucket:          get("STORAGE_BUCKET"),
		AccessKeyID:     get("STORAGE_ACCESS_KEY_ID"),
		SecretAccessKey: get("STORAGE_SECRET_ACCESS_KEY"),
		PublicBaseURL:   strings.TrimRight(get("STORAGE_PUBLIC_BASE_URL"), "/"),
	}
}

// Enabled reports whether storage was configured at all. An unset provider
// disables uploads without stopping the API.
func (c Config) Enabled() bool {
	return c.Provider != ""
}

// Validate rejects an unknown provider or a partially configured one, so a
// half-finished setup fails at startup instead of on the first upload.
func (c Config) Validate() error {
	required := map[string]string{
		"STORAGE_REGION":          c.Region,
		"STORAGE_BUCKET":          c.Bucket,
		"STORAGE_PUBLIC_BASE_URL": c.PublicBaseURL,
	}

	switch c.Provider {
	case ProviderSupabase:
		// Supabase has no ambient credentials or default endpoint.
		required["STORAGE_ENDPOINT"] = c.Endpoint
		required["STORAGE_ACCESS_KEY_ID"] = c.AccessKeyID
		required["STORAGE_SECRET_ACCESS_KEY"] = c.SecretAccessKey
	case ProviderS3:
		// Endpoint is optional (AWS default) and keys may come from the
		// default credential chain (e.g. an IAM role), but not half a pair.
		if (c.AccessKeyID == "") != (c.SecretAccessKey == "") {
			return fmt.Errorf("STORAGE_ACCESS_KEY_ID and STORAGE_SECRET_ACCESS_KEY must be set together")
		}
	default:
		return fmt.Errorf("unsupported STORAGE_PROVIDER %q (use %q or %q)", c.Provider, ProviderSupabase, ProviderS3)
	}

	var missing []string
	for _, name := range []string{"STORAGE_ENDPOINT", "STORAGE_REGION", "STORAGE_BUCKET", "STORAGE_ACCESS_KEY_ID", "STORAGE_SECRET_ACCESS_KEY", "STORAGE_PUBLIC_BASE_URL"} {
		if value, ok := required[name]; ok && value == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("STORAGE_PROVIDER=%s requires: %s", c.Provider, strings.Join(missing, ", "))
	}

	for name, raw := range map[string]string{"STORAGE_ENDPOINT": c.Endpoint, "STORAGE_PUBLIC_BASE_URL": c.PublicBaseURL} {
		if raw == "" {
			continue
		}
		if parsed, err := url.Parse(raw); err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.Host == "" {
			return fmt.Errorf("%s must be an absolute http(s) URL", name)
		}
	}

	return nil
}

// NewFromEnv builds the configured Storage. It returns (nil, nil) when
// STORAGE_PROVIDER is unset, and an error when the configuration is invalid.
func NewFromEnv(ctx context.Context) (Storage, error) {
	cfg := ConfigFromEnv()
	if !cfg.Enabled() {
		return nil, nil
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return newS3Storage(ctx, cfg)
}

// NewMenuImageKey returns a random, collision-free key for a menu image.
func NewMenuImageKey(extension string) (string, error) {
	return NewImageKey(MenuImagePrefix, extension)
}

// NewImageKey returns a random, collision-free key under prefix.
func NewImageKey(prefix, extension string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generating image key: %w", err)
	}
	return prefix + hex.EncodeToString(buf) + extension, nil
}

// IsManagedKey reports whether value is an object key this API stored, as
// opposed to an external URL (e.g. a placeholder) it must never delete.
func IsManagedKey(value string) bool {
	if strings.Contains(value, "://") {
		return false
	}
	for _, prefix := range managedPrefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

// ResolveURL turns a stored image value into the URL clients load: managed
// keys get the public base URL, anything else (external URLs such as a
// placeholder) is returned as is. It is nil-safe for when storage is disabled.
func ResolveURL(s Storage, value string) string {
	if s == nil || !IsManagedKey(value) {
		return value
	}
	return s.PublicBaseURL() + "/" + value
}

// KeyFromPublicURL maps one of this bucket's public URLs back to its key, so
// a client that echoes a resolved URL still results in a key being stored.
func KeyFromPublicURL(s Storage, value string) string {
	if s == nil {
		return value
	}
	if key, ok := strings.CutPrefix(value, s.PublicBaseURL()+"/"); ok && IsManagedKey(key) {
		return key
	}
	return value
}
