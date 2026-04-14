package services

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type fakePromotionStore struct {
	mu         sync.Mutex
	byID       map[int]*models.Promotion
	byCode     map[string]*models.Promotion
	usageCount map[int]int
}

func newFakePromotionStore(promotions ...*models.Promotion) *fakePromotionStore {
	store := &fakePromotionStore{
		byID:       make(map[int]*models.Promotion),
		byCode:     make(map[string]*models.Promotion),
		usageCount: make(map[int]int),
	}
	for _, p := range promotions {
		copyP := *p
		store.byID[p.ID] = &copyP
		store.byCode[p.Code] = &copyP
	}
	return store
}

func (f *fakePromotionStore) Create(promotion *models.Promotion) (*models.Promotion, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	copyP := *promotion
	if copyP.ID == 0 {
		copyP.ID = len(f.byID) + 1
	}
	f.byID[copyP.ID] = &copyP
	f.byCode[copyP.Code] = &copyP
	return &copyP, nil
}

func (f *fakePromotionStore) GetByID(id int) (*models.Promotion, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p := f.byID[id]
	if p == nil {
		return nil, nil
	}
	copyP := *p
	return &copyP, nil
}

func (f *fakePromotionStore) GetByCode(code string) (*models.Promotion, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p := f.byCode[code]
	if p == nil {
		return nil, nil
	}
	copyP := *p
	return &copyP, nil
}

func (f *fakePromotionStore) GetAll(limit, offset int, status string) ([]models.Promotion, int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	items := make([]models.Promotion, 0, len(f.byID))
	for _, p := range f.byID {
		items = append(items, *p)
	}
	return items, len(items), nil
}

func (f *fakePromotionStore) Update(promotion *models.Promotion) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	copyP := *promotion
	f.byID[promotion.ID] = &copyP
	f.byCode[promotion.Code] = &copyP
	return nil
}

func (f *fakePromotionStore) Delete(id int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	p := f.byID[id]
	if p != nil {
		delete(f.byCode, p.Code)
	}
	delete(f.byID, id)
	return nil
}

func (f *fakePromotionStore) ConsumeUsage(usage *models.PromotionUsage) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	p := f.byID[usage.PromotionID]
	if p == nil {
		return fmt.Errorf("promotion not found")
	}

	now := time.Now()
	if !p.IsActive || now.Before(p.ValidFrom) || now.After(p.ValidUntil) {
		return fmt.Errorf("promotion usage limit reached")
	}
	if p.UsageLimit != nil && p.UsedCount >= *p.UsageLimit {
		return fmt.Errorf("promotion usage limit reached")
	}

	p.UsedCount++
	f.usageCount[p.ID]++
	return nil
}

func TestValidatePromotionPercentageDiscountRespectsMaxCap(t *testing.T) {
	now := time.Now()
	maxDiscount := 15.0
	store := newFakePromotionStore(&models.Promotion{
		ID:            10,
		Code:          "PERCENT30",
		PromotionType: models.PromotionTypePercentage,
		DiscountValue: 30,
		MaxDiscount:   &maxDiscount,
		ValidFrom:     now.Add(-time.Hour),
		ValidUntil:    now.Add(time.Hour),
		IsActive:      true,
	})

	svc := NewPromotionService(store)

	validation, err := svc.ValidatePromotion("PERCENT30", 100, nil)
	if err != nil {
		t.Fatalf("ValidatePromotion() error = %v", err)
	}
	if !validation.IsValid {
		t.Fatalf("expected promotion to be valid, got message: %s", validation.Message)
	}
	if validation.Discount != 15 {
		t.Fatalf("expected capped discount 15, got %v", validation.Discount)
	}
}

func TestValidatePromotionFixedDiscountDoesNotExceedOrderAmount(t *testing.T) {
	now := time.Now()
	store := newFakePromotionStore(&models.Promotion{
		ID:            11,
		Code:          "FIXED50",
		PromotionType: models.PromotionTypeFixed,
		DiscountValue: 50,
		ValidFrom:     now.Add(-time.Hour),
		ValidUntil:    now.Add(time.Hour),
		IsActive:      true,
	})

	svc := NewPromotionService(store)

	validation, err := svc.ValidatePromotion("FIXED50", 20, nil)
	if err != nil {
		t.Fatalf("ValidatePromotion() error = %v", err)
	}
	if !validation.IsValid {
		t.Fatalf("expected promotion to be valid, got message: %s", validation.Message)
	}
	if validation.Discount != 20 {
		t.Fatalf("expected fixed discount to be capped at order amount (20), got %v", validation.Discount)
	}
}

func TestApplyPromotionConcurrentUsageLimitOne(t *testing.T) {
	now := time.Now()
	limit := 1
	store := newFakePromotionStore(&models.Promotion{
		ID:            12,
		Code:          "ONEONLY",
		PromotionType: models.PromotionTypeFixed,
		DiscountValue: 5,
		UsageLimit:    &limit,
		UsedCount:     0,
		ValidFrom:     now.Add(-time.Hour),
		ValidUntil:    now.Add(time.Hour),
		IsActive:      true,
	})

	svc := NewPromotionService(store)

	var wg sync.WaitGroup
	wg.Add(2)

	results := make(chan error, 2)
	for i := 0; i < 2; i++ {
		orderID := 100 + i
		go func() {
			defer wg.Done()
			_, err := svc.ApplyPromotionByCode("ONEONLY", orderID, nil, 25)
			results <- err
		}()
	}
	wg.Wait()
	close(results)

	successes := 0
	failures := 0
	for err := range results {
		if err == nil {
			successes++
		} else {
			failures++
		}
	}

	if successes != 1 || failures != 1 {
		t.Fatalf("expected 1 success and 1 failure, got %d success / %d failure", successes, failures)
	}
}
