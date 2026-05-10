package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMenuRepositorySearch_NoFilters_CountQueryDoesNotUsePaginationArgs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer db.Close()

	repo := NewMenuRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM menu_items WHERE is_available = true")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	createdAt := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, category_id, name, description, price, image_url, is_available, is_pre_order, preparation_time, tags, nutritional_info, created_at FROM menu_items WHERE is_available = true ORDER BY category_id, name LIMIT $1")).
		WithArgs(20).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"category_id",
				"name",
				"description",
				"price",
				"image_url",
				"is_available",
				"is_pre_order",
				"preparation_time",
				"tags",
				"nutritional_info",
				"created_at",
			}).AddRow(
				1,
				1,
				"Test Item",
				nil,
				9.99,
				nil,
				true,
				false,
				nil,
				[]byte("[]"),
				nil,
				createdAt,
			),
		)

	items, total, err := repo.Search("", nil, 20, 0)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if total != 1 {
		t.Fatalf("Search() total = %d, want 1", total)
	}
	if len(items) != 1 {
		t.Fatalf("Search() len(items) = %d, want 1", len(items))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
