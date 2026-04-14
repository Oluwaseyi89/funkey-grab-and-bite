package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type stubAdminLookup struct {
	admins map[int]*models.AdminUser
	err    error
}

func (s stubAdminLookup) GetAdminUserByID(adminID int) (*models.AdminUser, error) {
	if s.err != nil {
		return nil, s.err
	}

	admin := s.admins[adminID]
	if admin == nil {
		return nil, nil
	}

	adminCopy := *admin
	return &adminCopy, nil
}

func TestAdminAuthMiddlewareAllowsActiveAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	originalSecret := string(utils.GetSecret())
	utils.SetSecret("test-admin-secret")
	t.Cleanup(func() {
		utils.SetSecret(originalSecret)
	})

	token, err := utils.GenerateAdminJWT(42, "admin@funkey.com", "admin")
	if err != nil {
		t.Fatalf("GenerateAdminJWT() error = %v", err)
	}

	lookup := stubAdminLookup{
		admins: map[int]*models.AdminUser{
			42: {ID: 42, Email: "admin@funkey.com", IsActive: true},
		},
	}

	recorder := httptest.NewRecorder()
	router := gin.New()
	router.Use(AdminAuthMiddleware(lookup))
	router.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestAdminAuthMiddlewareRejectsForgedAdminWithoutAdminRecord(t *testing.T) {
	gin.SetMode(gin.TestMode)
	originalSecret := string(utils.GetSecret())
	utils.SetSecret("test-admin-secret")
	t.Cleanup(func() {
		utils.SetSecret(originalSecret)
	})

	token, err := utils.GenerateAdminJWT(404, "forged@funkey.com", "admin")
	if err != nil {
		t.Fatalf("GenerateAdminJWT() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	router := gin.New()
	router.Use(AdminAuthMiddleware(stubAdminLookup{}))
	router.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}
}

func TestAdminAuthMiddlewareRejectsDisabledAdminWithExistingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	originalSecret := string(utils.GetSecret())
	utils.SetSecret("test-admin-secret")
	t.Cleanup(func() {
		utils.SetSecret(originalSecret)
	})

	token, err := utils.GenerateAdminJWT(77, "disabled@funkey.com", "admin")
	if err != nil {
		t.Fatalf("GenerateAdminJWT() error = %v", err)
	}

	lookup := stubAdminLookup{
		admins: map[int]*models.AdminUser{
			77: {ID: 77, Email: "disabled@funkey.com", IsActive: false},
		},
	}

	recorder := httptest.NewRecorder()
	router := gin.New()
	router.Use(AdminAuthMiddleware(lookup))
	router.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}
}

func TestAdminAuthMiddlewareRevokesExistingTokenAfterAdminDeactivation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	originalSecret := string(utils.GetSecret())
	utils.SetSecret("test-admin-secret")
	t.Cleanup(func() {
		utils.SetSecret(originalSecret)
	})

	token, err := utils.GenerateAdminJWT(91, "revoked@funkey.com", "admin")
	if err != nil {
		t.Fatalf("GenerateAdminJWT() error = %v", err)
	}

	lookup := stubAdminLookup{
		admins: map[int]*models.AdminUser{
			91: {ID: 91, Email: "revoked@funkey.com", IsActive: true},
		},
	}

	router := gin.New()
	router.Use(AdminAuthMiddleware(lookup))
	router.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	first := httptest.NewRecorder()
	firstReq := httptest.NewRequest(http.MethodGet, "/admin", nil)
	firstReq.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(first, firstReq)
	if first.Code != http.StatusOK {
		t.Fatalf("expected first request status %d, got %d", http.StatusOK, first.Code)
	}

	lookup.admins[91].IsActive = false

	second := httptest.NewRecorder()
	secondReq := httptest.NewRequest(http.MethodGet, "/admin", nil)
	secondReq.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(second, secondReq)
	if second.Code != http.StatusForbidden {
		t.Fatalf("expected second request status %d after deactivation, got %d", http.StatusForbidden, second.Code)
	}
}
