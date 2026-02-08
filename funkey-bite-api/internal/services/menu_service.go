package services

import (
	"fmt"
	"strings"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

type MenuService interface {
	GetMenuItems() ([]models.MenuItem, error)
	GetMenuItemByID(id int) (*models.MenuItem, error)
	GetMenuItemsByCategory(categoryID int) ([]models.MenuItem, error)
	GetCategories() ([]models.MenuCategory, error)
	SearchMenuItems(query string, categoryID *int, page, limit int) ([]models.MenuItem, int, error)
	GetFeaturedItems(limit int) ([]models.MenuItem, error)
	GetMenuItemsByTags(tags []string) ([]models.MenuItem, error)
	GetMenuWithFilters(filters MenuFilters) ([]models.MenuItem, int, error)
}

type MenuFilters struct {
	Query      string
	CategoryID *int
	Tags       []string
	MinPrice   *float64
	MaxPrice   *float64
	IsPreOrder *bool
	Page       int
	Limit      int
}

type menuService struct {
	menuRepo repository.MenuRepository
}

func NewMenuService(menuRepo repository.MenuRepository) MenuService {
	return &menuService{
		menuRepo: menuRepo,
	}
}

func (s *menuService) GetMenuItems() ([]models.MenuItem, error) {
	items, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get menu items: %w", err)
	}
	return items, nil
}

func (s *menuService) GetMenuItemByID(id int) (*models.MenuItem, error) {
	item, err := s.menuRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}
	if item == nil {
		return nil, fmt.Errorf("menu item not found")
	}
	return item, nil
}

func (s *menuService) GetMenuItemsByCategory(categoryID int) ([]models.MenuItem, error) {
	items, err := s.menuRepo.GetByCategory(categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu items by category: %w", err)
	}
	return items, nil
}

func (s *menuService) GetCategories() ([]models.MenuCategory, error) {
	categories, err := s.menuRepo.GetCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

func (s *menuService) SearchMenuItems(query string, categoryID *int, page, limit int) ([]models.MenuItem, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	items, total, err := s.menuRepo.Search(query, categoryID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search menu items: %w", err)
	}

	return items, total, nil
}

func (s *menuService) GetFeaturedItems(limit int) ([]models.MenuItem, error) {
	if limit < 1 {
		limit = 6
	}
	if limit > 20 {
		limit = 20
	}

	return s.menuRepo.GetFeaturedItems(limit)
}

func (s *menuService) GetMenuItemsByTags(tags []string) ([]models.MenuItem, error) {
	if len(tags) == 0 {
		return s.GetMenuItems()
	}

	return s.menuRepo.GetByTags(tags)
}

func (s *menuService) GetMenuWithFilters(filters MenuFilters) ([]models.MenuItem, int, error) {
	return s.SearchMenuItems(filters.Query, filters.CategoryID, filters.Page, filters.Limit)
}

// Keep for backward compatibility - maintains the old SearchMenuItems signature
func (s *menuService) SearchMenuItemsOld(query string) ([]models.MenuItem, error) {
	allItems, err := s.GetMenuItems()
	if err != nil {
		return nil, err
	}

	var results []models.MenuItem
	query = strings.ToLower(query)

	for _, item := range allItems {
		if strings.Contains(strings.ToLower(item.Name), query) ||
			strings.Contains(strings.ToLower(item.Description), query) {
			results = append(results, item)
		}
	}

	return results, nil
}

// package services

// import (
// 	"fmt"
// 	"strings"

// 	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
// 	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
// )

// type MenuService interface {
// 	GetMenuItems() ([]models.MenuItem, error)
// 	GetMenuItemByID(id int) (*models.MenuItem, error)
// 	GetMenuItemsByCategory(categoryID int) ([]models.MenuItem, error)
// 	GetCategories() ([]models.MenuCategory, error)
// 	SearchMenuItems(query string) ([]models.MenuItem, error)
// }

// type menuService struct {
// 	menuRepo repository.MenuRepository
// }

// func NewMenuService(menuRepo repository.MenuRepository) MenuService {
// 	return &menuService{
// 		menuRepo: menuRepo,
// 	}
// }

// func (s *menuService) GetMenuItems() ([]models.MenuItem, error) {
// 	items, err := s.menuRepo.GetAll()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get menu items: %w", err)
// 	}
// 	return items, nil
// }

// func (s *menuService) GetMenuItemByID(id int) (*models.MenuItem, error) {
// 	item, err := s.menuRepo.GetByID(id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get menu item: %w", err)
// 	}
// 	if item == nil {
// 		return nil, fmt.Errorf("menu item not found")
// 	}
// 	return item, nil
// }

// func (s *menuService) GetMenuItemsByCategory(categoryID int) ([]models.MenuItem, error) {
// 	items, err := s.menuRepo.GetByCategory(categoryID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get menu items by category: %w", err)
// 	}
// 	return items, nil
// }

// func (s *menuService) GetCategories() ([]models.MenuCategory, error) {
// 	categories, err := s.menuRepo.GetCategories()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get categories: %w", err)
// 	}
// 	return categories, nil
// }

// func (s *menuService) SearchMenuItems(query string) ([]models.MenuItem, error) {
// 	// This would require adding a search method to the repository
// 	// For now, we'll filter from all items
// 	allItems, err := s.GetMenuItems()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var results []models.MenuItem
// 	query = strings.ToLower(query)

// 	for _, item := range allItems {
// 		if strings.Contains(strings.ToLower(item.Name), query) ||
// 			strings.Contains(strings.ToLower(item.Description), query) {
// 			results = append(results, item)
// 		}
// 	}

// 	return results, nil
// }
