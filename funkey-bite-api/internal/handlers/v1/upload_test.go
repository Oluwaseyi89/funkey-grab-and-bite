package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeImageStorage struct {
	putKeys      []string
	contentTypes []string
}

func (f *fakeImageStorage) Put(_ context.Context, key, contentType string, body io.Reader, _ int64) error {
	_, _ = io.Copy(io.Discard, body)
	f.putKeys = append(f.putKeys, key)
	f.contentTypes = append(f.contentTypes, contentType)
	return nil
}
func (f *fakeImageStorage) Delete(context.Context, string) error { return nil }
func (f *fakeImageStorage) PublicBaseURL() string                { return "https://cdn.example.com/assets" }

// pngBytes is a valid PNG header padded to the requested size.
func pngBytes(size int) []byte {
	header := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")
	return append(header, bytes.Repeat([]byte{0}, size-len(header))...)
}

func uploadRequest(t *testing.T, field, filename string, content []byte) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(field, filename)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write(content)
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func serveUpload(handler *UploadHandler, req *http.Request) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/upload", handler.UploadMenuImage)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func TestUploadMenuImageStoresSniffedImage(t *testing.T) {
	store := &fakeImageStorage{}
	// The client claims a .jpg, but the bytes are a PNG: the sniffed type wins.
	recorder := serveUpload(NewUploadHandler(store), uploadRequest(t, "file", "photo.jpg", pngBytes(1024)))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}
	var resp struct{ Key, URL string }
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(resp.Key, "menu-items/") || !strings.HasSuffix(resp.Key, ".png") {
		t.Fatalf("unexpected key %q", resp.Key)
	}
	if resp.URL != "https://cdn.example.com/assets/"+resp.Key {
		t.Fatalf("unexpected url %q", resp.URL)
	}
	if len(store.putKeys) != 1 || store.putKeys[0] != resp.Key || store.contentTypes[0] != "image/png" {
		t.Fatalf("unexpected storage calls: keys=%v types=%v", store.putKeys, store.contentTypes)
	}
}

func TestUploadMenuImageRejections(t *testing.T) {
	cases := []struct {
		name    string
		handler *UploadHandler
		req     func(t *testing.T) *http.Request
		want    int
	}{
		{"storage not configured", NewUploadHandler(nil), func(t *testing.T) *http.Request {
			return uploadRequest(t, "file", "a.png", pngBytes(1024))
		}, http.StatusServiceUnavailable},
		{"missing file field", NewUploadHandler(&fakeImageStorage{}), func(t *testing.T) *http.Request {
			return uploadRequest(t, "image", "a.png", pngBytes(1024))
		}, http.StatusBadRequest},
		{"not an image", NewUploadHandler(&fakeImageStorage{}), func(t *testing.T) *http.Request {
			return uploadRequest(t, "file", "a.png", []byte("<html><script>alert(1)</script></html>"))
		}, http.StatusUnsupportedMediaType},
		{"too large", NewUploadHandler(&fakeImageStorage{}), func(t *testing.T) *http.Request {
			return uploadRequest(t, "file", "a.png", pngBytes(MaxMenuImageBytes+1))
		}, http.StatusRequestEntityTooLarge},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := serveUpload(tc.handler, tc.req(t))
			if recorder.Code != tc.want {
				t.Fatalf("expected %d, got %d: %s", tc.want, recorder.Code, recorder.Body.String())
			}
		})
	}
}
