package v1

import (
	"context"
	"net/http"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type totalsCall struct{ from, to string }

type fakeReportsStore struct {
	totals []totalsCall
	series []totalsCall
	group  string
}

func (f *fakeReportsStore) Totals(_ context.Context, _ string, from, to time.Time) (models.PeriodTotals, error) {
	f.totals = append(f.totals, totalsCall{from.Format("2006-01-02 15:04"), to.Format("2006-01-02 15:04")})
	return models.PeriodTotals{}, nil
}
func (f *fakeReportsStore) Series(_ context.Context, _ string, from, to time.Time, group string) ([]models.SalesPoint, error) {
	f.series = append(f.series, totalsCall{from.Format("2006-01-02"), to.Format("2006-01-02")})
	f.group = group
	return []models.SalesPoint{}, nil
}
func (f *fakeReportsStore) Hourly(context.Context, string, time.Time, time.Time) ([]models.HourPoint, error) {
	return nil, nil
}
func (f *fakeReportsStore) TopItems(context.Context, string, time.Time, time.Time, int) ([]models.TopItem, error) {
	return nil, nil
}
func (f *fakeReportsStore) PaymentMethods(context.Context, string, time.Time, time.Time) ([]models.Share, error) {
	return nil, nil
}
func (f *fakeReportsStore) OrderTypes(context.Context, string, time.Time, time.Time) ([]models.Share, error) {
	return nil, nil
}
func (f *fakeReportsStore) Workload(context.Context) (models.OrderQueue, int, int, int, int, error) {
	return models.OrderQueue{}, 0, 0, 0, 0, nil
}

func reportsRouter(store *fakeReportsStore, now time.Time) *gin.Engine {
	gin.SetMode(gin.TestMode)
	lagos, _ := time.LoadLocation("Africa/Lagos")
	h := NewReportsHandler(store, lagos)
	h.now = func() time.Time { return now }
	r := gin.New()
	r.GET("/overview", h.GetOverview)
	r.GET("/summary", h.GetSalesSummary)
	return r
}

func TestSalesSummaryIncludesTheLastDayAndComparesTheDaysBefore(t *testing.T) {
	store := &fakeReportsStore{}
	rec := call(reportsRouter(store, time.Now()), http.MethodGet, "/summary?from=2026-11-01&to=2026-11-07&group=week", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if store.totals[0] != (totalsCall{"2026-11-01 00:00", "2026-11-08 00:00"}) {
		t.Fatalf("the period must run to the end of the last day, got %+v", store.totals[0])
	}
	if store.totals[1] != (totalsCall{"2026-10-25 00:00", "2026-11-01 00:00"}) {
		t.Fatalf("the comparison must be the 7 days before, got %+v", store.totals[1])
	}
	if store.group != "week" {
		t.Fatalf("expected weekly grouping passed on, got %q", store.group)
	}
}

func TestSalesSummaryRejectsBadInput(t *testing.T) {
	for _, query := range []string{
		"from=2026-11-07&to=2026-11-01",
		"from=yesterday&to=2026-11-01",
		"from=2026-11-01&to=2026-11-02&group=year",
		"from=2020-01-01&to=2026-11-01",
	} {
		if rec := call(reportsRouter(&fakeReportsStore{}, time.Now()), http.MethodGet, "/summary?"+query, ""); rec.Code != http.StatusBadRequest {
			t.Fatalf("%s: expected 400, got %d", query, rec.Code)
		}
	}
}

func TestOverviewUsesTheBusinessDayAndComparesYesterdaySoFar(t *testing.T) {
	store := &fakeReportsStore{}
	// 23:30 UTC on 10 Nov is 00:30 on 11 Nov in Lagos.
	now := time.Date(2026, 11, 10, 23, 30, 0, 0, time.UTC)
	rec := call(reportsRouter(store, now), http.MethodGet, "/overview", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if store.totals[0] != (totalsCall{"2026-11-11 00:00", "2026-11-11 00:30"}) {
		t.Fatalf("today must start at midnight in Lagos, got %+v", store.totals[0])
	}
	if store.totals[1] != (totalsCall{"2026-11-10 00:00", "2026-11-10 00:30"}) {
		t.Fatalf("yesterday must stop at the same time of day, got %+v", store.totals[1])
	}
	if store.series[0] != (totalsCall{"2026-10-29", "2026-11-11"}) {
		t.Fatalf("expected the last 14 days including today, got %+v", store.series[0])
	}
}
