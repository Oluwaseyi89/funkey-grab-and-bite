package storage

import (
	"context"
	"io"
	"regexp"
	"strings"
	"testing"
)

type stubStorage struct{ base string }

func (s stubStorage) Put(context.Context, string, string, io.Reader, int64) error { return nil }
func (s stubStorage) Delete(context.Context, string) error                        { return nil }
func (s stubStorage) PublicBaseURL() string                                       { return s.base }

func validSupabaseConfig() Config {
	return Config{
		Provider:        ProviderSupabase,
		Endpoint:        "https://ref.supabase.co/storage/v1/s3",
		Region:          "us-west-2",
		Bucket:          "assets",
		AccessKeyID:     "id",
		SecretAccessKey: "secret",
		PublicBaseURL:   "https://ref.supabase.co/storage/v1/object/public/assets",
	}
}

func TestConfigValidate(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
	}{
		{"valid supabase", func(*Config) {}, ""},
		{"supabase needs endpoint", func(c *Config) { c.Endpoint = "" }, "STORAGE_ENDPOINT"},
		{"supabase needs keys", func(c *Config) { c.AccessKeyID, c.SecretAccessKey = "", "" }, "STORAGE_ACCESS_KEY_ID, STORAGE_SECRET_ACCESS_KEY"},
		{"s3 without endpoint or keys", func(c *Config) {
			c.Provider, c.Endpoint, c.AccessKeyID, c.SecretAccessKey = ProviderS3, "", "", ""
		}, ""},
		{"s3 rejects half a key pair", func(c *Config) { c.Provider, c.SecretAccessKey = ProviderS3, "" }, "must be set together"},
		{"needs bucket and base url", func(c *Config) { c.Bucket, c.PublicBaseURL = "", "" }, "STORAGE_BUCKET, STORAGE_PUBLIC_BASE_URL"},
		{"unknown provider", func(c *Config) { c.Provider = "gcs" }, "unsupported STORAGE_PROVIDER"},
		{"base url must be absolute", func(c *Config) { c.PublicBaseURL = "assets/public" }, "STORAGE_PUBLIC_BASE_URL must be an absolute"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validSupabaseConfig()
			tc.mutate(&cfg)
			err := cfg.Validate()
			if tc.wantErr == "" && err != nil {
				t.Fatalf("expected valid config, got %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestNewFromEnvDisabledWhenProviderUnset(t *testing.T) {
	t.Setenv("STORAGE_PROVIDER", "")
	s, err := NewFromEnv(context.Background())
	if s != nil || err != nil {
		t.Fatalf("expected (nil, nil) when STORAGE_PROVIDER is unset, got (%v, %v)", s, err)
	}
}

func TestNewMenuImageKey(t *testing.T) {
	key, err := NewMenuImageKey(".webp")
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^menu-items/[0-9a-f]{32}\.webp$`).MatchString(key) {
		t.Fatalf("unexpected key format %q", key)
	}
	if other, _ := NewMenuImageKey(".webp"); other == key {
		t.Fatal("expected distinct keys")
	}
}

func TestResolveURLAndKeyFromPublicURL(t *testing.T) {
	s := stubStorage{base: "https://cdn.example.com/assets"}
	key := "menu-items/abc.png"
	placeholder := "https://via.placeholder.com/400x300"

	if got := ResolveURL(s, key); got != "https://cdn.example.com/assets/menu-items/abc.png" {
		t.Fatalf("ResolveURL(key) = %q", got)
	}
	if got := ResolveURL(s, placeholder); got != placeholder {
		t.Fatalf("external URLs must pass through, got %q", got)
	}
	if got := ResolveURL(nil, key); got != key {
		t.Fatalf("nil storage must pass values through, got %q", got)
	}

	if got := KeyFromPublicURL(s, ResolveURL(s, key)); got != key {
		t.Fatalf("KeyFromPublicURL round trip = %q", got)
	}
	if got := KeyFromPublicURL(s, placeholder); got != placeholder {
		t.Fatalf("foreign URLs must be kept, got %q", got)
	}
	if got := KeyFromPublicURL(s, "https://cdn.example.com/assets/other/x.png"); got != "https://cdn.example.com/assets/other/x.png" {
		t.Fatalf("only menu image keys are extracted, got %q", got)
	}
}

func TestIsManagedKey(t *testing.T) {
	for value, want := range map[string]bool{
		"menu-items/a.jpg":                   true,
		"":                                   false,
		"https://x.test/menu-items/a.jpg":    false,
		"other/a.jpg":                        false,
		"menu-items/https://evil.test/a.jpg": false,
	} {
		if got := IsManagedKey(value); got != want {
			t.Errorf("IsManagedKey(%q) = %t, want %t", value, got, want)
		}
	}
}
