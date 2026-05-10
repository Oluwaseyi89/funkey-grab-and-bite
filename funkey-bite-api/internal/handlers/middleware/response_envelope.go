package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"funkey-grab-and-bite/funkey-bite-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

// responseCaptureWriter buffers JSON responses so we can normalize shape once handlers finish.
type responseCaptureWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *responseCaptureWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (w *responseCaptureWriter) WriteHeaderNow() {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
}

func (w *responseCaptureWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	return w.body.Write(data)
}

func (w *responseCaptureWriter) WriteString(s string) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	return w.body.WriteString(s)
}

func (w *responseCaptureWriter) Status() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}

	return w.statusCode
}

func (w *responseCaptureWriter) Size() int {
	return w.body.Len()
}

func (w *responseCaptureWriter) Written() bool {
	return w.statusCode != 0 || w.body.Len() > 0
}

func ResponseEnvelopeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		capture := &responseCaptureWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
			statusCode:     http.StatusOK,
		}
		c.Writer = capture

		c.Next()

		status := capture.Status()
		rawBody := bytes.TrimSpace(capture.body.Bytes())
		if len(rawBody) == 0 {
			capture.ResponseWriter.WriteHeader(status)
			return
		}

		contentType := capture.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			capture.ResponseWriter.WriteHeader(status)
			_, _ = capture.ResponseWriter.Write(rawBody)
			return
		}

		var payload interface{}
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			capture.ResponseWriter.WriteHeader(status)
			_, _ = capture.ResponseWriter.Write(rawBody)
			return
		}

		if isAPIEnvelope(payload) {
			capture.ResponseWriter.WriteHeader(status)
			_, _ = capture.ResponseWriter.Write(rawBody)
			return
		}

		wrapped := handlers.APIResponse{Success: status < http.StatusBadRequest}
		if status >= http.StatusBadRequest {
			code, message := deriveErrorFromPayload(payload, status)
			wrapped.Error = &handlers.APIError{
				Code:    code,
				Message: message,
			}
		} else {
			wrapped.Data = payload
		}

		normalized, err := json.Marshal(wrapped)
		if err != nil {
			capture.ResponseWriter.WriteHeader(status)
			_, _ = capture.ResponseWriter.Write(rawBody)
			return
		}

		capture.Header().Set("Content-Type", "application/json; charset=utf-8")
		capture.ResponseWriter.WriteHeader(status)
		_, _ = capture.ResponseWriter.Write(normalized)
	}
}

func isAPIEnvelope(payload interface{}) bool {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return false
	}

	_, hasSuccess := data["success"]
	_, hasData := data["data"]
	_, hasError := data["error"]
	return hasSuccess && (hasData || hasError)
}

func deriveErrorFromPayload(payload interface{}, status int) (string, string) {
	defaultCode := strings.ToUpper(strings.ReplaceAll(http.StatusText(status), " ", "_"))
	defaultMessage := http.StatusText(status)

	obj, ok := payload.(map[string]interface{})
	if !ok {
		return defaultCode, defaultMessage
	}

	if errorValue, ok := obj["error"]; ok {
		switch typed := errorValue.(type) {
		case string:
			if typed != "" {
				return defaultCode, typed
			}
		case map[string]interface{}:
			code, _ := typed["code"].(string)
			message, _ := typed["message"].(string)
			if code == "" {
				code = defaultCode
			}
			if message == "" {
				message = defaultMessage
			}
			return code, message
		}
	}

	if message, ok := obj["message"].(string); ok && message != "" {
		return defaultCode, message
	}

	return defaultCode, defaultMessage
}
