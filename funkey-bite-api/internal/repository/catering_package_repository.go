package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type CateringPackageRepository struct {
	db *sql.DB
}

func NewCateringPackageRepository(db *sql.DB) *CateringPackageRepository {
	return &CateringPackageRepository{db: db}
}

const cateringPackageColumns = `
	id, name, description, price_per_person, min_guests, max_guests,
	service_style, min_notice_days, inclusions, dish_choices, is_active,
	display_order, created_at, updated_at
`

// List returns packages in display order; activeOnly hides switched-off ones.
func (r *CateringPackageRepository) List(activeOnly bool) ([]models.CateringPackage, error) {
	query := `SELECT ` + cateringPackageColumns + ` FROM catering_packages
		WHERE ($1 = false OR is_active = true)
		ORDER BY display_order, id`

	rows, err := r.db.Query(query, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to list catering packages: %w", err)
	}
	defer rows.Close()

	packages := []models.CateringPackage{}
	for rows.Next() {
		pkg, err := scanCateringPackage(rows)
		if err != nil {
			return nil, err
		}
		packages = append(packages, *pkg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read catering packages: %w", err)
	}

	if err := r.attachItems(packages); err != nil {
		return nil, err
	}
	return packages, nil
}

func (r *CateringPackageRepository) GetByID(id int) (*models.CateringPackage, error) {
	row := r.db.QueryRow(`SELECT `+cateringPackageColumns+` FROM catering_packages WHERE id = $1`, id)
	pkg, err := scanCateringPackage(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	packages := []models.CateringPackage{*pkg}
	if err := r.attachItems(packages); err != nil {
		return nil, err
	}
	return &packages[0], nil
}

// Create inserts a package and its menu items in one transaction.
func (r *CateringPackageRepository) Create(input *models.CateringPackageInput) (*models.CateringPackage, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	inclusions, err := json.Marshal(nonNilStrings(input.Inclusions))
	if err != nil {
		return nil, fmt.Errorf("failed to encode inclusions: %w", err)
	}

	var id int
	err = tx.QueryRow(`
		INSERT INTO catering_packages (
			name, description, price_per_person, min_guests, max_guests,
			service_style, min_notice_days, inclusions, dish_choices,
			is_active, display_order
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`,
		input.Name, input.Description, input.PricePerPerson, input.MinGuests, input.MaxGuests,
		input.ServiceStyle, input.MinNoticeDays, inclusions, input.DishChoices,
		input.IsActive, input.DisplayOrder,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create catering package: %w", err)
	}

	if err := replacePackageItems(tx, id, input.MenuItemIDs); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit catering package: %w", err)
	}
	return r.GetByID(id)
}

// Update replaces a package and its menu items; it returns nil if there is no
// package with that id.
func (r *CateringPackageRepository) Update(id int, input *models.CateringPackageInput) (*models.CateringPackage, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	inclusions, err := json.Marshal(nonNilStrings(input.Inclusions))
	if err != nil {
		return nil, fmt.Errorf("failed to encode inclusions: %w", err)
	}

	result, err := tx.Exec(`
		UPDATE catering_packages SET
			name = $1, description = $2, price_per_person = $3, min_guests = $4,
			max_guests = $5, service_style = $6, min_notice_days = $7,
			inclusions = $8, dish_choices = $9, is_active = $10,
			display_order = $11, updated_at = $12
		WHERE id = $13`,
		input.Name, input.Description, input.PricePerPerson, input.MinGuests,
		input.MaxGuests, input.ServiceStyle, input.MinNoticeDays,
		inclusions, input.DishChoices, input.IsActive,
		input.DisplayOrder, time.Now(), id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update catering package: %w", err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, nil
	}

	if err := replacePackageItems(tx, id, input.MenuItemIDs); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit catering package: %w", err)
	}
	return r.GetByID(id)
}

func replacePackageItems(tx *sql.Tx, packageID int, menuItemIDs []int) error {
	if _, err := tx.Exec(`DELETE FROM catering_package_items WHERE package_id = $1`, packageID); err != nil {
		return fmt.Errorf("failed to clear package menu items: %w", err)
	}
	for _, menuItemID := range menuItemIDs {
		if _, err := tx.Exec(
			`INSERT INTO catering_package_items (package_id, menu_item_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			packageID, menuItemID,
		); err != nil {
			return fmt.Errorf("failed to link menu item %d to package: %w", menuItemID, err)
		}
	}
	return nil
}

// attachItems loads the linked menu items for the given packages.
func (r *CateringPackageRepository) attachItems(packages []models.CateringPackage) error {
	if len(packages) == 0 {
		return nil
	}

	index := make(map[int]*models.CateringPackage, len(packages))
	for i := range packages {
		packages[i].MenuItems = []models.CateringPackageItem{}
		index[packages[i].ID] = &packages[i]
	}

	rows, err := r.db.Query(`
		SELECT cpi.package_id, mi.id, mi.name, mi.is_available
		FROM catering_package_items cpi
		JOIN menu_items mi ON mi.id = cpi.menu_item_id
		ORDER BY mi.name`)
	if err != nil {
		return fmt.Errorf("failed to load package menu items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var packageID int
		var item models.CateringPackageItem
		if err := rows.Scan(&packageID, &item.MenuItemID, &item.Name, &item.IsAvailable); err != nil {
			return fmt.Errorf("failed to scan package menu item: %w", err)
		}
		if pkg, ok := index[packageID]; ok {
			pkg.MenuItems = append(pkg.MenuItems, item)
		}
	}
	return rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanCateringPackage(row rowScanner) (*models.CateringPackage, error) {
	var pkg models.CateringPackage
	var maxGuests sql.NullInt64
	var inclusions []byte

	err := row.Scan(
		&pkg.ID, &pkg.Name, &pkg.Description, &pkg.PricePerPerson, &pkg.MinGuests, &maxGuests,
		&pkg.ServiceStyle, &pkg.MinNoticeDays, &inclusions, &pkg.DishChoices, &pkg.IsActive,
		&pkg.DisplayOrder, &pkg.CreatedAt, &pkg.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan catering package: %w", err)
	}

	if maxGuests.Valid {
		max := int(maxGuests.Int64)
		pkg.MaxGuests = &max
	}
	pkg.Inclusions = []string{}
	if len(inclusions) > 0 {
		if err := json.Unmarshal(inclusions, &pkg.Inclusions); err != nil {
			return nil, fmt.Errorf("failed to parse package inclusions: %w", err)
		}
	}
	return &pkg, nil
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
