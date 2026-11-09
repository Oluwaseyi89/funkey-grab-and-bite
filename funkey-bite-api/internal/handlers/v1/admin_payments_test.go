package v1

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

type fakeOnlinePayments struct{ filter models.OnlinePaymentFilter }

func (f *fakeOnlinePayments) List(ctx context.Context, filter models.OnlinePaymentFilter, limit, offset int) ([]models.OnlinePayment, int, error) {
	f.filter = filter
	return []models.OnlinePayment{{OrderID: 7, OrderNumber: "FG-7", PaymentStatus: "needs_review"}}, 1, nil
}
func (f *fakeOnlinePayments) CountByStatus(ctx context.Context) (map[string]int, error) {
	return map[string]int{"needs_review": 1}, nil
}
func (f *fakeOnlinePayments) Unmatched(ctx context.Context, limit int) ([]models.PaymentEventRecord, error) {
	return []models.PaymentEventRecord{}, nil
}

type fakeReviewer struct {
	err      error
	decision string
	note     string
}

func (f *fakeReviewer) ResolvePaymentReview(ctx context.Context, orderID int, decision, note string) (string, error) {
	f.decision, f.note = decision, note
	return "paid", f.err
}

type fakeJobAdmin struct{ retried bool }

func (f *fakeJobAdmin) List(ctx context.Context, status string, limit, offset int) ([]models.JobSummary, int, error) {
	return []models.JobSummary{}, 0, nil
}
func (f *fakeJobAdmin) RetryDead(ctx context.Context, id int64) (bool, error) { return f.retried, nil }

type fakeAuditReader struct{ filter models.AuditFilter }

func (f *fakeAuditReader) List(ctx context.Context, filter models.AuditFilter, limit, offset int) ([]models.AuditRecord, int, error) {
	f.filter = filter
	return []models.AuditRecord{}, 0, nil
}

func adminPaymentsRouter(h *AdminPaymentsHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/admin/payments", h.ListPayments)
	router.POST("/admin/orders/:id/payment-review", h.ResolvePaymentReview)
	router.POST("/admin/jobs/:id/retry", h.RetryJob)
	router.GET("/admin/audit", h.ListAudit)
	router.GET("/admin/jobs", h.ListJobs)
	return router
}

func call(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestListPaymentsPassesFiltersAndCounts(t *testing.T) {
	payments := &fakeOnlinePayments{}
	router := adminPaymentsRouter(NewAdminPaymentsHandler(payments, &fakeReviewer{}, &fakePaymentService{}, &fakeJobAdmin{}, &fakeAuditReader{}))

	rec := call(router, http.MethodGet, "/admin/payments?status=needs_review&search=%20FG-7%20", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"needs_review":1`) {
		t.Fatalf("expected payments with counts, got %d %s", rec.Code, rec.Body.String())
	}
	if payments.filter.PaymentStatus != "needs_review" || payments.filter.Search != "FG-7" {
		t.Fatalf("expected the filters passed on trimmed, got %+v", payments.filter)
	}
}

func TestResolvePaymentReview(t *testing.T) {
	cases := map[string]struct {
		body     string
		err      error
		wantCode int
	}{
		"accept":           {`{"decision":"accept","note":"Customer sent the right amount later"}`, nil, http.StatusOK},
		"unknown decision": {`{"decision":"ignore"}`, nil, http.StatusBadRequest},
		"not under review": {`{"decision":"accept"}`, repository.ErrNotUnderReview, http.StatusConflict},
		"cannot reopen":    {`{"decision":"accept"}`, repository.ErrCannotReopen, http.StatusConflict},
		"order not found":  {`{"decision":"refund"}`, sql.ErrNoRows, http.StatusNotFound},
		"note too long":    {`{"decision":"refund","note":"` + strings.Repeat("x", 501) + `"}`, nil, http.StatusBadRequest},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			reviewer := &fakeReviewer{err: tc.err}
			router := adminPaymentsRouter(NewAdminPaymentsHandler(&fakeOnlinePayments{}, reviewer, &fakePaymentService{}, &fakeJobAdmin{}, &fakeAuditReader{}))
			rec := call(router, http.MethodPost, "/admin/orders/7/payment-review", tc.body)
			if rec.Code != tc.wantCode {
				t.Fatalf("expected %d, got %d: %s", tc.wantCode, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestRetryJobOnlyRetriesFailedJobs(t *testing.T) {
	router := adminPaymentsRouter(NewAdminPaymentsHandler(&fakeOnlinePayments{}, &fakeReviewer{}, &fakePaymentService{}, &fakeJobAdmin{retried: false}, &fakeAuditReader{}))
	if rec := call(router, http.MethodPost, "/admin/jobs/5/retry", ""); rec.Code != http.StatusConflict {
		t.Fatalf("a job that isn't dead must not be retried, got %d", rec.Code)
	}
	if rec := call(router, http.MethodGet, "/admin/jobs?status=everything", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("an unknown status filter must be rejected, got %d", rec.Code)
	}
}

func TestListAuditPassesFilters(t *testing.T) {
	audit := &fakeAuditReader{}
	router := adminPaymentsRouter(NewAdminPaymentsHandler(&fakeOnlinePayments{}, &fakeReviewer{}, &fakePaymentService{}, &fakeJobAdmin{}, audit))
	call(router, http.MethodGet, "/admin/audit?entityType=order&entityId=7&actorType=paystack", "")
	if audit.filter != (models.AuditFilter{EntityType: "order", EntityID: "7", ActorType: "paystack"}) {
		t.Fatalf("expected the filters passed on, got %+v", audit.filter)
	}
}
