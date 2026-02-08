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

// Add these methods to MenuRepository:

// Search searches for menu items using database full-text search
// func (r *MenuRepository) Search(query string, categoryID *int, limit, offset int) ([]models.MenuItem, int, error) {
// 	// Build the search query
// 	baseQuery := `
//         SELECT id, category_id, name, description, price, image_url,
//                is_available, is_pre_order, preparation_time, tags,
//                nutritional_info, created_at
//         FROM menu_items
//         WHERE is_available = true
//     `

// 	countQuery := `SELECT COUNT(*) FROM menu_items WHERE is_available = true`

// 	// Build WHERE conditions
// 	var conditions []string
// 	var params []interface{}
// 	paramIndex := 1

// 	// Add search condition
// 	if query != "" {
// 		// Using PostgreSQL full-text search with tsvector
// 		// First, create a searchable column in your database:
// 		// ALTER TABLE menu_items ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (
// 		//     setweight(to_tsvector('english', COALESCE(name, '')), 'A') ||
// 		//     setweight(to_tsvector('english', COALESCE(description, '')), 'B') ||
// 		//     setweight(to_tsvector('english', COALESCE(array_to_string(tags, ' '), '')), 'C')
// 		// ) STORED;
// 		// CREATE INDEX idx_menu_items_search ON menu_items USING GIN(search_vector);

// 		conditions = append(conditions,
// 			`(search_vector @@ plainto_tsquery('english', $`+fmt.Sprint(paramIndex)+`) OR
//               name ILIKE $`+fmt.Sprint(paramIndex+1)+` OR
//               description ILIKE $`+fmt.Sprint(paramIndex+2)+`)`)
// 		searchQuery := "%" + query + "%"
// 		params = append(params, query, searchQuery, searchQuery)
// 		paramIndex += 3
// 	}

// 	// Add category filter
// 	if categoryID != nil {
// 		conditions = append(conditions, `category_id = $`+fmt.Sprint(paramIndex))
// 		params = append(params, *categoryID)
// 		paramIndex++
// 	}

// 	// Apply conditions
// 	if len(conditions) > 0 {
// 		whereClause := " AND " + strings.Join(conditions, " AND ")
// 		baseQuery += whereClause
// 		countQuery += whereClause
// 	}

// 	// Add ordering (by search relevance if query provided, otherwise by name)
// 	if query != "" {
// 		baseQuery += ` ORDER BY
//             ts_rank(search_vector, plainto_tsquery('english', $1)) DESC,
//             name ASC`
// 	} else {
// 		baseQuery += ` ORDER BY category_id, name`
// 	}

// 	// Add pagination
// 	if limit > 0 {
// 		baseQuery += fmt.Sprintf(` LIMIT $%d`, paramIndex)
// 		params = append(params, limit)
// 		paramIndex++

// 		if offset > 0 {
// 			baseQuery += fmt.Sprintf(` OFFSET $%d`, paramIndex)
// 			params = append(params, offset)
// 		}
// 	}

// 	// Get total count
// 	var total int
// 	countParams := params
// 	if query != "" {
// 		// Remove LIMIT/OFFSET params for count query
// 		if limit > 0 {
// 			countParams = countParams[:len(countParams)-1]
// 			if offset > 0 {
// 				countParams = countParams[:len(countParams)-1]
// 			}
// 		}
// 	}

// 	err := r.db.QueryRow(countQuery, countParams...).Scan(&total)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get search count: %w", err)
// 	}

// 	// Get search results
// 	rows, err := r.db.Query(baseQuery, params...)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to search menu items: %w", err)
// 	}
// 	defer rows.Close()

// 	return r.scanMenuItems(rows)
// }

func (r *MenuRepository) Search(query string, categoryID *int, limit, offset int) ([]models.MenuItem, int, error) {
	// Build the search query
	baseQuery := `
        SELECT id, category_id, name, description, price, image_url, 
               is_available, is_pre_order, preparation_time, tags, 
               nutritional_info, created_at
        FROM menu_items
        WHERE is_available = true
    `

	countQuery := `SELECT COUNT(*) FROM menu_items WHERE is_available = true`

	// Build WHERE conditions
	var conditions []string
	var params []interface{}
	paramIndex := 1

	// Add search condition
	if query != "" {
		conditions = append(conditions,
			`(search_vector @@ plainto_tsquery('english', $`+fmt.Sprint(paramIndex)+`) OR 
              name ILIKE $`+fmt.Sprint(paramIndex+1)+` OR 
              description ILIKE $`+fmt.Sprint(paramIndex+2)+`)`)
		searchQuery := "%" + query + "%"
		params = append(params, query, searchQuery, searchQuery)
		paramIndex += 3
	}

	// Add category filter
	if categoryID != nil {
		conditions = append(conditions, `category_id = $`+fmt.Sprint(paramIndex))
		params = append(params, *categoryID)
		paramIndex++
	}

	// Apply conditions
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Add ordering (by search relevance if query provided, otherwise by name)
	if query != "" {
		baseQuery += ` ORDER BY 
            ts_rank(search_vector, plainto_tsquery('english', $1)) DESC,
            name ASC`
	} else {
		baseQuery += ` ORDER BY category_id, name`
	}

	// Add pagination
	if limit > 0 {
		baseQuery += fmt.Sprintf(` LIMIT $%d`, paramIndex)
		params = append(params, limit)
		paramIndex++

		if offset > 0 {
			baseQuery += fmt.Sprintf(` OFFSET $%d`, paramIndex)
			params = append(params, offset)
		}
	}

	// Get total count
	var total int
	countParams := params
	if query != "" {
		// Remove LIMIT/OFFSET params for count query
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

	// Get search results
	rows, err := r.db.Query(baseQuery, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search menu items: %w", err)
	}
	defer rows.Close()

	// Scan the rows into menu items
	items, err := r.scanMenuItems(rows)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// GetFeaturedItems gets featured/popular menu items
func (r *MenuRepository) GetFeaturedItems(limit int) ([]models.MenuItem, error) {
	// This query gets popular items based on order history
	// You might want to cache this since it's computationally expensive

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

// GetByTags gets menu items by tags
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

// Helper method to scan menu items (extract from your existing GetAll method)
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

		// Parse tags JSON
		var tags []string
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
