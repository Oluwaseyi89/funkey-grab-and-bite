package repository

import (
	"database/sql"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type PromotionRepository struct {
	db *sql.DB
}

func NewPromotionRepository(db *sql.DB) *PromotionRepository {
	return &PromotionRepository{db: db}
}

func (r *PromotionRepository) Create(promotion *models.Promotion) (*models.Promotion, error) {
	query := `
        INSERT INTO promotions (
            code, title, description, promotion_type, discount_value, max_discount,
            min_order_amount, valid_from, valid_until, usage_limit, used_count,
            is_active, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        RETURNING id, created_at, updated_at
    `

	now := time.Now()
	err := r.db.QueryRow(
		query,
		promotion.Code,
		promotion.Title,
		promotion.Description,
		promotion.PromotionType,
		promotion.DiscountValue,
		promotion.MaxDiscount,
		promotion.MinOrderAmount,
		promotion.ValidFrom,
		promotion.ValidUntil,
		promotion.UsageLimit,
		promotion.UsedCount,
		promotion.IsActive,
		now,
		now,
	).Scan(&promotion.ID, &promotion.CreatedAt, &promotion.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create promotion: %w", err)
	}

	return promotion, nil
}

func (r *PromotionRepository) GetByID(id int) (*models.Promotion, error) {
	query := `
        SELECT id, code, title, description, promotion_type, discount_value, max_discount,
               min_order_amount, valid_from, valid_until, usage_limit, used_count,
               is_active, created_at, updated_at
        FROM promotions
        WHERE id = $1
    `

	var promotion models.Promotion
	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&promotion.ID,
		&promotion.Code,
		&promotion.Title,
		&promotion.Description,
		&promotion.PromotionType,
		&promotion.DiscountValue,
		&promotion.MaxDiscount,
		&promotion.MinOrderAmount,
		&promotion.ValidFrom,
		&promotion.ValidUntil,
		&promotion.UsageLimit,
		&promotion.UsedCount,
		&promotion.IsActive,
		&promotion.CreatedAt,
		&promotion.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get promotion: %w", err)
	}

	return &promotion, nil
}

func (r *PromotionRepository) GetByCode(code string) (*models.Promotion, error) {
	query := `
        SELECT id, code, title, description, promotion_type, discount_value, max_discount,
               min_order_amount, valid_from, valid_until, usage_limit, used_count,
               is_active, created_at, updated_at
        FROM promotions
        WHERE code = $1
    `

	var promotion models.Promotion
	row := r.db.QueryRow(query, code)

	err := row.Scan(
		&promotion.ID,
		&promotion.Code,
		&promotion.Title,
		&promotion.Description,
		&promotion.PromotionType,
		&promotion.DiscountValue,
		&promotion.MaxDiscount,
		&promotion.MinOrderAmount,
		&promotion.ValidFrom,
		&promotion.ValidUntil,
		&promotion.UsageLimit,
		&promotion.UsedCount,
		&promotion.IsActive,
		&promotion.CreatedAt,
		&promotion.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get promotion: %w", err)
	}

	return &promotion, nil
}

func (r *PromotionRepository) GetAll(limit, offset int, status string) ([]models.Promotion, int, error) {
	countQuery := `SELECT COUNT(*) FROM promotions WHERE ($1 = '' OR is_active = $1::boolean)`
	var total int
	err := r.db.QueryRow(countQuery, status).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get promotion count: %w", err)
	}

	query := `
        SELECT id, code, title, description, promotion_type, discount_value, max_discount,
               min_order_amount, valid_from, valid_until, usage_limit, used_count,
               is_active, created_at, updated_at
        FROM promotions
        WHERE ($1 = '' OR is_active = $1::boolean)
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get promotions: %w", err)
	}
	defer rows.Close()

	var promotions []models.Promotion
	for rows.Next() {
		var promotion models.Promotion
		err := rows.Scan(
			&promotion.ID,
			&promotion.Code,
			&promotion.Title,
			&promotion.Description,
			&promotion.PromotionType,
			&promotion.DiscountValue,
			&promotion.MaxDiscount,
			&promotion.MinOrderAmount,
			&promotion.ValidFrom,
			&promotion.ValidUntil,
			&promotion.UsageLimit,
			&promotion.UsedCount,
			&promotion.IsActive,
			&promotion.CreatedAt,
			&promotion.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan promotion: %w", err)
		}
		promotions = append(promotions, promotion)
	}

	return promotions, total, nil
}

func (r *PromotionRepository) Update(promotion *models.Promotion) error {
	query := `
        UPDATE promotions SET
            title = $1,
            description = $2,
            promotion_type = $3,
            discount_value = $4,
            max_discount = $5,
            min_order_amount = $6,
            valid_from = $7,
            valid_until = $8,
            usage_limit = $9,
            is_active = $10,
            updated_at = $11
        WHERE id = $12
    `

	_, err := r.db.Exec(
		query,
		promotion.Title,
		promotion.Description,
		promotion.PromotionType,
		promotion.DiscountValue,
		promotion.MaxDiscount,
		promotion.MinOrderAmount,
		promotion.ValidFrom,
		promotion.ValidUntil,
		promotion.UsageLimit,
		promotion.IsActive,
		time.Now(),
		promotion.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update promotion: %w", err)
	}

	return nil
}

func (r *PromotionRepository) Delete(id int) error {
	query := `DELETE FROM promotions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PromotionRepository) IncrementUsage(id int) error {
	query := `UPDATE promotions SET used_count = used_count + 1, updated_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

func (r *PromotionRepository) RecordUsage(usage *models.PromotionUsage) error {
	query := `
        INSERT INTO promotion_usage (promotion_id, order_id, customer_id, discount_applied, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	err := r.db.QueryRow(
		query,
		usage.PromotionID,
		usage.OrderID,
		usage.CustomerID,
		usage.DiscountApplied,
		time.Now(),
	).Scan(&usage.ID)

	if err != nil {
		return fmt.Errorf("failed to record promotion usage: %w", err)
	}

	return nil
}
