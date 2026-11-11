package v1

import (
	"context"
	"log"
	"net/http"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

// reportsStore is what the reports need from the database.
type reportsStore interface {
	Totals(ctx context.Context, tz string, from, to time.Time) (models.PeriodTotals, error)
	Series(ctx context.Context, tz string, from, to time.Time, group string) ([]models.SalesPoint, error)
	Hourly(ctx context.Context, tz string, from, to time.Time) ([]models.HourPoint, error)
	TopItems(ctx context.Context, tz string, from, to time.Time, limit int) ([]models.TopItem, error)
	PaymentMethods(ctx context.Context, tz string, from, to time.Time) ([]models.Share, error)
	OrderTypes(ctx context.Context, tz string, from, to time.Time) ([]models.Share, error)
	Workload(ctx context.Context) (models.OrderQueue, int, int, int, int, error)
}

// ReportsHandler serves the admin dashboard and sales reports, with days
// taken in the business timezone.
type ReportsHandler struct {
	store    reportsStore
	location *time.Location
	now      func() time.Time
}

func NewReportsHandler(store reportsStore, location *time.Location) *ReportsHandler {
	return &ReportsHandler{store: store, location: location, now: time.Now}
}

// The longest period a report may cover, so one request can't scan years.
const maxReportDays = 731

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetOverview returns today so far, what needs doing, and the last 14 days.
// @Summary Dashboard overview
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.DashboardOverview
// @Router /admin/dashboard/overview [get]
func (h *ReportsHandler) GetOverview(c *gin.Context) {
	ctx := c.Request.Context()
	tz := h.location.String()
	now := h.now().In(h.location)
	today := startOfDay(now)
	tomorrow := today.AddDate(0, 0, 1)

	overview := models.DashboardOverview{Timezone: tz}
	var err error
	fail := func(what string, err error) {
		log.Printf("dashboard overview: %s: %v", what, err)
		handlers.Error(c, http.StatusInternalServerError, "REPORT_FAILED", "Couldn't load the dashboard")
	}

	if overview.Today, err = h.store.Totals(ctx, tz, today, now); err != nil {
		fail("today", err)
		return
	}
	// Yesterday up to this time of day, so the morning isn't compared with a whole day.
	if overview.YesterdaySoFar, err = h.store.Totals(ctx, tz, today.AddDate(0, 0, -1), now.AddDate(0, 0, -1)); err != nil {
		fail("yesterday", err)
		return
	}
	if overview.Daily, err = h.store.Series(ctx, tz, today.AddDate(0, 0, -13), today, "day"); err != nil {
		fail("daily", err)
		return
	}
	if overview.Hourly, err = h.store.Hourly(ctx, tz, today, tomorrow); err != nil {
		fail("hourly", err)
		return
	}
	if overview.TopItems, err = h.store.TopItems(ctx, tz, today.AddDate(0, 0, -6), tomorrow, 5); err != nil {
		fail("top items", err)
		return
	}
	overview.Queue, overview.PaymentsToReview, overview.AwaitingPayment, overview.LowStock, overview.CateringToQuote, err = h.store.Workload(ctx)
	if err != nil {
		fail("workload", err)
		return
	}

	handlers.Success(c, overview)
}

// GetSalesSummary returns the sales for a period (from and to are dates,
// both included), grouped by day, week or month, with the previous period of
// the same length for comparison.
// @Summary Sales summary
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param from query string true "First day, YYYY-MM-DD"
// @Param to query string true "Last day, YYYY-MM-DD"
// @Param group query string false "day, week or month"
// @Success 200 {object} models.SalesSummary
// @Router /admin/reports/summary [get]
func (h *ReportsHandler) GetSalesSummary(c *gin.Context) {
	from, errFrom := time.ParseInLocation("2006-01-02", c.Query("from"), h.location)
	to, errTo := time.ParseInLocation("2006-01-02", c.Query("to"), h.location)
	if errFrom != nil || errTo != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_DATES", "Give from and to as dates, e.g. 2026-11-01")
		return
	}
	if from.After(to) {
		handlers.Error(c, http.StatusBadRequest, "INVALID_DATES", "The start date must be on or before the end date")
		return
	}
	days := int(to.Sub(from).Hours()/24+0.5) + 1
	if days > maxReportDays {
		handlers.Error(c, http.StatusBadRequest, "RANGE_TOO_LONG", "Choose a period of up to two years")
		return
	}
	group := c.DefaultQuery("group", "day")
	if group != "day" && group != "week" && group != "month" {
		handlers.Error(c, http.StatusBadRequest, "INVALID_GROUP", "Group by day, week or month")
		return
	}

	ctx := c.Request.Context()
	tz := h.location.String()
	end := to.AddDate(0, 0, 1)
	prevTo := from.AddDate(0, 0, -1)
	prevFrom := from.AddDate(0, 0, -days)

	summary := models.SalesSummary{
		From: from.Format("2006-01-02"), To: to.Format("2006-01-02"), Group: group, Timezone: tz,
		PreviousFrom: prevFrom.Format("2006-01-02"), PreviousTo: prevTo.Format("2006-01-02"),
	}
	var err error
	fail := func(what string, err error) {
		log.Printf("sales summary: %s: %v", what, err)
		handlers.Error(c, http.StatusInternalServerError, "REPORT_FAILED", "Couldn't load the report")
	}
	if summary.Totals, err = h.store.Totals(ctx, tz, from, end); err != nil {
		fail("totals", err)
		return
	}
	if summary.Previous, err = h.store.Totals(ctx, tz, prevFrom, from); err != nil {
		fail("previous", err)
		return
	}
	if summary.Series, err = h.store.Series(ctx, tz, from, to, group); err != nil {
		fail("series", err)
		return
	}
	if summary.Hourly, err = h.store.Hourly(ctx, tz, from, end); err != nil {
		fail("hourly", err)
		return
	}
	if summary.TopItems, err = h.store.TopItems(ctx, tz, from, end, 10); err != nil {
		fail("top items", err)
		return
	}
	if summary.PaymentMethods, err = h.store.PaymentMethods(ctx, tz, from, end); err != nil {
		fail("payment methods", err)
		return
	}
	if summary.OrderTypes, err = h.store.OrderTypes(ctx, tz, from, end); err != nil {
		fail("order types", err)
		return
	}

	handlers.Success(c, summary)
}
