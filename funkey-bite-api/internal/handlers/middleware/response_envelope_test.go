package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResponseEnvelopeMiddleware_WrapsLegacySuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ResponseEnvelopeMiddleware())
	r.GET("/legacy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"value": 42})
	})

	req := httptest.NewRequest(http.MethodGet, "/legacy", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if ok, _ := body["success"].(bool); !ok {
		t.Fatalf("expected success=true, got %v", body["success"])
	}

	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected object data payload, got %T", body["data"])
	}

	if data["value"] != float64(42) {
		t.Fatalf("expected data.value=42, got %v", data["value"])
	}
}

func TestResponseEnvelopeMiddleware_WrapsLegacyError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ResponseEnvelopeMiddleware())
	r.GET("/legacy-error", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
	})

	req := httptest.NewRequest(http.MethodGet, "/legacy-error", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if ok, _ := body["success"].(bool); ok {
		t.Fatalf("expected success=false, got %v", body["success"])
	}

	errorObj, ok := body["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected error object, got %T", body["error"])
	}

	if errorObj["message"] != "bad input" {
		t.Fatalf("expected error.message=bad input, got %v", errorObj["message"])
	}
}

func TestResponseEnvelopeMiddleware_LeavesStandardEnvelopeUntouched(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ResponseEnvelopeMiddleware())
	r.GET("/standard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"ok": true}})
	})

	req := httptest.NewRequest(http.MethodGet, "/standard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if ok, _ := body["success"].(bool); !ok {
		t.Fatalf("expected success=true, got %v", body["success"])
	}

	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object, got %T", body["data"])
	}

	if data["ok"] != true {
		t.Fatalf("expected data.ok=true, got %v", data["ok"])
	}
}
