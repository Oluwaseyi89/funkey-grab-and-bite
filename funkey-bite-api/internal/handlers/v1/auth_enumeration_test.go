package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type fakeAuthCheckService struct {
	user   *models.User
	exists bool
	err    error
}

func (f *fakeAuthCheckService) Register(userData models.UserRegistration) (*models.AuthResponse, error) {
	return nil, nil
}
func (f *fakeAuthCheckService) Login(phone, password string) (*models.AuthResponse, error) {
	return nil, nil
}
func (f *fakeAuthCheckService) CheckUserExists(phone, email string) (*models.User, bool, error) {
	return f.user, f.exists, f.err
}
func (f *fakeAuthCheckService) AuthenticateOrder(orderData models.OrderWithAuth) (*models.User, error) {
	return nil, nil
}

type fakeUserProfileService struct{}

func (f *fakeUserProfileService) GetByID(id int) (*models.User, error) { return nil, nil }
func (f *fakeUserProfileService) UpdateProfile(id int, updates *models.ProfileUpdate) (*models.User, error) {
	return nil, nil
}
func (f *fakeUserProfileService) GetOrderHistory(userID int) ([]models.Order, error) { return nil, nil }
func (f *fakeUserProfileService) ChangePassword(userID int, currentPassword, newPassword string) error {
	return nil
}

func TestCheckUserRequiresPhoneOrEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAuthHandler(&fakeAuthCheckService{}, &fakeUserProfileService{})
	router := gin.New()
	router.GET("/auth/check", handler.CheckUser)

	req := httptest.NewRequest(http.MethodGet, "/auth/check", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d body=%s", http.StatusBadRequest, res.Code, res.Body.String())
	}
}

func TestCheckUserDoesNotExposeExistsOrUserForExistingAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	email := "exists@example.com"
	handler := NewAuthHandler(&fakeAuthCheckService{
		user:   &models.User{ID: 10, Phone: "+15551230000", Email: &email, FullName: "Existing User"},
		exists: true,
	}, &fakeUserProfileService{})

	router := gin.New()
	router.GET("/auth/check", handler.CheckUser)

	req := httptest.NewRequest(http.MethodGet, "/auth/check?phone=+15551230000", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d body=%s", http.StatusOK, res.Code, res.Body.String())
	}

	body := res.Body.String()
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}
	if _, ok := payload["exists"]; ok {
		t.Fatalf("auth check leaked exists field: %s", body)
	}
	if _, ok := payload["user"]; ok {
		t.Fatalf("auth check leaked user field: %s", body)
	}
	if !strings.Contains(body, "If an account exists for the provided details") {
		t.Fatalf("expected generic message, got body=%s", body)
	}
}

func TestCheckUserResponseIsSameForExistingAndMissingAccounts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	routerExisting := gin.New()
	routerExisting.GET("/auth/check", NewAuthHandler(&fakeAuthCheckService{exists: true}, &fakeUserProfileService{}).CheckUser)
	reqExisting := httptest.NewRequest(http.MethodGet, "/auth/check?email=exists@example.com", nil)
	resExisting := httptest.NewRecorder()
	routerExisting.ServeHTTP(resExisting, reqExisting)

	routerMissing := gin.New()
	routerMissing.GET("/auth/check", NewAuthHandler(&fakeAuthCheckService{exists: false}, &fakeUserProfileService{}).CheckUser)
	reqMissing := httptest.NewRequest(http.MethodGet, "/auth/check?email=missing@example.com", nil)
	resMissing := httptest.NewRecorder()
	routerMissing.ServeHTTP(resMissing, reqMissing)

	if resExisting.Code != http.StatusOK || resMissing.Code != http.StatusOK {
		t.Fatalf("expected 200 for both responses, got existing=%d missing=%d", resExisting.Code, resMissing.Code)
	}
	if resExisting.Body.String() != resMissing.Body.String() {
		t.Fatalf("auth check response differs by account existence and enables enumeration: existing=%s missing=%s", resExisting.Body.String(), resMissing.Body.String())
	}
}
