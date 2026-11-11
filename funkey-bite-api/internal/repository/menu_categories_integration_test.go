package repository

import (
	"fmt"
	"testing"
	"time"
)

// Hidden categories stay off the public menu but stay listed for the admin,
// so they can be shown again.
func TestHiddenCategoriesAreListedForTheAdminOnly(t *testing.T) {
	db := openIntegrationDB(t)
	name := fmt.Sprintf("Hidden test %d", time.Now().UnixNano())
	var id int
	if err := db.QueryRow(`INSERT INTO menu_categories (name, description, display_order, is_active) VALUES ($1, '', 999, false) RETURNING id`, name).Scan(&id); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Exec(`DELETE FROM menu_categories WHERE id = $1`, id) })

	repo := NewMenuRepository(db)
	has := func(list []string) bool {
		for _, n := range list {
			if n == name {
				return true
			}
		}
		return false
	}
	public, err := repo.GetCategories()
	if err != nil {
		t.Fatal(err)
	}
	all, err := repo.GetAllCategories()
	if err != nil {
		t.Fatal(err)
	}
	names := func(n int, get func(int) string) []string {
		out := make([]string, n)
		for i := range out {
			out[i] = get(i)
		}
		return out
	}
	if has(names(len(public), func(i int) string { return public[i].Name })) {
		t.Fatal("a hidden category must not be on the public menu")
	}
	if !has(names(len(all), func(i int) string { return all[i].Name })) {
		t.Fatal("the admin must still see a hidden category")
	}
}

// The admin's menu list includes dishes marked unavailable; the customer
// menu doesn't.
func TestUnavailableDishesAreListedForTheAdminOnly(t *testing.T) {
	db := openIntegrationDB(t)
	name := fmt.Sprintf("Unavailable test %d", time.Now().UnixNano())
	var id int
	if err := db.QueryRow(`INSERT INTO menu_items (category_id, name, description, price, is_available) VALUES (1, $1, '', 100, false) RETURNING id`, name).Scan(&id); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Exec(`DELETE FROM menu_items WHERE id = $1`, id) })

	repo := NewMenuRepository(db)
	public, _, err := repo.Search(name, nil, 50, 0)
	if err != nil {
		t.Fatal(err)
	}
	all, _, err := repo.SearchAll(name, nil, 50, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(public) != 0 {
		t.Fatal("an unavailable dish must not be on the customer menu")
	}
	if len(all) != 1 || all[0].ID != id {
		t.Fatalf("the admin must still find an unavailable dish, got %d", len(all))
	}
}
