package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type MenuRepository struct {
	db *sql.DB
}

func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) GetByID(id int) (*models.MenuItem, error) {
	query := `
		SELECT id, category_id, name, description, price, image_url, 
		       is_available, is_pre_order, preparation_time, tags, 
		       nutritional_info
		FROM menu_items
		WHERE id = $1
	`

	var item models.MenuItem
	var tagsJSON []byte
	var nutritionalInfoJSON []byte
	var tags []string

	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&item.ID,
		&item.CategoryID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.ImageURL,
		&item.IsAvailable,
		&item.IsPreOrder,
		&item.PreparationTime,
		&tagsJSON,
		&nutritionalInfoJSON,
		&item.CreatedAt, // Add this line

	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}

	// Parse tags JSON
	if len(tagsJSON) > 0 {
		if err := json.Unmarshal(tagsJSON, &tags); err != nil {
			return nil, fmt.Errorf("failed to parse tags: %w", err)
		}
		item.Tags = tags
	} else {
		item.Tags = []string{}
	}

	// Parse nutritional info JSON
	if len(nutritionalInfoJSON) > 0 {
		var nutritionalInfo models.NutritionalInfo
		if err := json.Unmarshal(nutritionalInfoJSON, &nutritionalInfo); err != nil {
			return nil, fmt.Errorf("failed to parse nutritional info: %w", err)
		}
		item.NutritionalInfo = &nutritionalInfo
	}

	return &item, nil
}

func (r *MenuRepository) GetAll() ([]models.MenuItem, error) {
	query := `
		SELECT id, category_id, name, description, price, image_url, 
		       is_available, is_pre_order, preparation_time, tags, 
		       nutritional_info
		FROM menu_items
		WHERE is_available = true
		ORDER BY category_id, name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu items: %w", err)
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		var tagsJSON []byte
		var nutritionalInfoJSON []byte
		var tags []string

		err := rows.Scan(
			&item.ID,
			&item.CategoryID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.ImageURL,
			&item.IsAvailable,
			&item.IsPreOrder,
			&item.PreparationTime,
			&tagsJSON,
			&nutritionalInfoJSON,
			&item.CreatedAt, // Add this line

		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu item: %w", err)
		}

		// Parse tags JSON
		if len(tagsJSON) > 0 {
			if err := json.Unmarshal(tagsJSON, &tags); err != nil {
				return nil, fmt.Errorf("failed to parse tags: %w", err)
			}
			item.Tags = tags
		} else {
			item.Tags = []string{}
		}

		// Parse nutritional info JSON
		if len(nutritionalInfoJSON) > 0 {
			var nutritionalInfo models.NutritionalInfo
			if err := json.Unmarshal(nutritionalInfoJSON, &nutritionalInfo); err != nil {
				return nil, fmt.Errorf("failed to parse nutritional info: %w", err)
			}
			item.NutritionalInfo = &nutritionalInfo
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *MenuRepository) GetByCategory(categoryID int) ([]models.MenuItem, error) {
	query := `
		SELECT id, category_id, name, description, price, image_url, 
		       is_available, is_pre_order, preparation_time, tags, 
		       nutritional_info
		FROM menu_items
		WHERE category_id = $1 AND is_available = true
		ORDER BY name
	`

	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu items by category: %w", err)
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		var tagsJSON []byte
		var nutritionalInfoJSON []byte
		var tags []string

		err := rows.Scan(
			&item.ID,
			&item.CategoryID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.ImageURL,
			&item.IsAvailable,
			&item.IsPreOrder,
			&item.PreparationTime,
			&tagsJSON,
			&nutritionalInfoJSON,
			&item.CreatedAt, // Add this line

		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu item: %w", err)
		}

		// Parse tags JSON
		if len(tagsJSON) > 0 {
			if err := json.Unmarshal(tagsJSON, &tags); err != nil {
				return nil, fmt.Errorf("failed to parse tags: %w", err)
			}
			item.Tags = tags
		} else {
			item.Tags = []string{}
		}

		// Parse nutritional info JSON
		if len(nutritionalInfoJSON) > 0 {
			var nutritionalInfo models.NutritionalInfo
			if err := json.Unmarshal(nutritionalInfoJSON, &nutritionalInfo); err != nil {
				return nil, fmt.Errorf("failed to parse nutritional info: %w", err)
			}
			item.NutritionalInfo = &nutritionalInfo
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *MenuRepository) GetCategories() ([]models.MenuCategory, error) {
	query := `
		SELECT id, name, description, display_order, is_active
		FROM menu_categories
		WHERE is_active = true
		ORDER BY display_order
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu categories: %w", err)
	}
	defer rows.Close()

	var categories []models.MenuCategory
	for rows.Next() {
		var category models.MenuCategory

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.DisplayOrder,
			&category.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu category: %w", err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}
