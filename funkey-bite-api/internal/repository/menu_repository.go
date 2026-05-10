package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

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
		&item.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}

	if len(tagsJSON) > 0 {
		if err := json.Unmarshal(tagsJSON, &tags); err != nil {
			return nil, fmt.Errorf("failed to parse tags: %w", err)
		}
		item.Tags = tags
	} else {
		item.Tags = []string{}
	}

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
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu item: %w", err)
		}

		if len(tagsJSON) > 0 {
			if err := json.Unmarshal(tagsJSON, &tags); err != nil {
				return nil, fmt.Errorf("failed to parse tags: %w", err)
			}
			item.Tags = tags
		} else {
			item.Tags = []string{}
		}

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
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu item: %w", err)
		}

		if len(tagsJSON) > 0 {
			if err := json.Unmarshal(tagsJSON, &tags); err != nil {
				return nil, fmt.Errorf("failed to parse tags: %w", err)
			}
			item.Tags = tags
		} else {
			item.Tags = []string{}
		}

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
		SELECT id, name, description, display_order, is_active, created_at
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
			&category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu category: %w", err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (r *MenuRepository) GetCategoryByID(id int) (*models.MenuCategory, error) {
	query := `
		SELECT id, name, description, display_order, is_active, created_at
		FROM menu_categories
		WHERE id = $1
	`

	var category models.MenuCategory
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.DisplayOrder,
		&category.IsActive,
		&category.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get menu category: %w", err)
	}

	return &category, nil
}

func (r *MenuRepository) CreateCategory(category *models.MenuCategory) (*models.MenuCategory, error) {
	query := `
		INSERT INTO menu_categories (name, description, display_order, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		category.Name,
		category.Description,
		category.DisplayOrder,
		category.IsActive,
	).Scan(&category.ID, &category.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create menu category: %w", err)
	}

	return category, nil
}

func (r *MenuRepository) UpdateCategory(category *models.MenuCategory) error {
	query := `
		UPDATE menu_categories
		SET name = $1,
			description = $2,
			display_order = $3,
			is_active = $4
		WHERE id = $5
	`

	result, err := r.db.Exec(
		query,
		category.Name,
		category.Description,
		category.DisplayOrder,
		category.IsActive,
		category.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update menu category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check category update result: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *MenuRepository) Search(query string, categoryID *int, limit, offset int) ([]models.MenuItem, int, error) {
	baseQuery := `
        SELECT id, category_id, name, description, price, image_url, 
               is_available, is_pre_order, preparation_time, tags, 
               nutritional_info, created_at
        FROM menu_items
        WHERE is_available = true
    `

	countQuery := `SELECT COUNT(*) FROM menu_items WHERE is_available = true`

	var conditions []string
	var params []interface{}
	paramIndex := 1

	if query != "" {
		conditions = append(conditions,
			`(search_vector @@ plainto_tsquery('english', $`+fmt.Sprint(paramIndex)+`) OR 
              name ILIKE $`+fmt.Sprint(paramIndex+1)+` OR 
              description ILIKE $`+fmt.Sprint(paramIndex+2)+`)`)
		searchQuery := "%" + query + "%"
		params = append(params, query, searchQuery, searchQuery)
		paramIndex += 3
	}

	if categoryID != nil {
		conditions = append(conditions, `category_id = $`+fmt.Sprint(paramIndex))
		params = append(params, *categoryID)
		paramIndex++
	}

	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	if query != "" {
		baseQuery += ` ORDER BY 
            ts_rank(search_vector, plainto_tsquery('english', $1)) DESC,
            name ASC`
	} else {
		baseQuery += ` ORDER BY category_id, name`
	}

	if limit > 0 {
		baseQuery += fmt.Sprintf(` LIMIT $%d`, paramIndex)
		params = append(params, limit)
		paramIndex++

		if offset > 0 {
			baseQuery += fmt.Sprintf(` OFFSET $%d`, paramIndex)
			params = append(params, offset)
		}
	}

	var total int
	countParams := params
	if query != "" {
		if limit > 0 {
			countParams = countParams[:len(countParams)-1]
			if offset > 0 {
				countParams = countParams[:len(countParams)-1]
			}
		}
	}

	err := r.db.QueryRow(countQuery, countParams...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get search count: %w", err)
	}

	rows, err := r.db.Query(baseQuery, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search menu items: %w", err)
	}
	defer rows.Close()

	items, err := r.scanMenuItems(rows)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *MenuRepository) GetFeaturedItems(limit int) ([]models.MenuItem, error) {

	query := `
        SELECT mi.id, mi.category_id, mi.name, mi.description, mi.price, 
               mi.image_url, mi.is_available, mi.is_pre_order, 
               mi.preparation_time, mi.tags, mi.nutritional_info, mi.created_at
        FROM menu_items mi
        LEFT JOIN (
            SELECT menu_item_id, SUM(quantity) as total_ordered
            FROM order_items oi
            JOIN orders o ON oi.order_id = o.id
            WHERE o.created_at >= NOW() - INTERVAL '30 days'
            GROUP BY menu_item_id
        ) popular ON mi.id = popular.menu_item_id
        WHERE mi.is_available = true
        ORDER BY COALESCE(popular.total_ordered, 0) DESC, mi.name
        LIMIT $1
    `

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get featured items: %w", err)
	}
	defer rows.Close()

	return r.scanMenuItems(rows)
}

func (r *MenuRepository) GetByTags(tags []string) ([]models.MenuItem, error) {
	query := `
        SELECT id, category_id, name, description, price, image_url, 
               is_available, is_pre_order, preparation_time, tags, 
               nutritional_info, created_at
        FROM menu_items
        WHERE is_available = true AND tags @> $1
        ORDER BY name
    `

	tagsJSON, _ := json.Marshal(tags)
	rows, err := r.db.Query(query, tagsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu items by tags: %w", err)
	}
	defer rows.Close()

	return r.scanMenuItems(rows)
}

func (r *MenuRepository) scanMenuItems(rows *sql.Rows) ([]models.MenuItem, error) {
	var items []models.MenuItem

	for rows.Next() {
		var item models.MenuItem
		var tagsJSON []byte
		var nutritionalInfoJSON []byte

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
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu item: %w", err)
		}

		var tags []string
		if len(tagsJSON) > 0 {
			if err := json.Unmarshal(tagsJSON, &tags); err != nil {
				return nil, fmt.Errorf("failed to parse tags: %w", err)
			}
			item.Tags = tags
		} else {
			item.Tags = []string{}
		}

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
