package repository

import (
	"database/sql"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, phone, email, full_name, password_hash, is_verified, 
		       is_active, last_login, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
	`

	row := r.db.QueryRow(query, id)
	return r.scanUser(row)
}

func (r *UserRepository) FindByPhone(phone string) (*models.User, error) {
	query := `
		SELECT id, phone, email, full_name, password_hash, is_verified, 
		       is_active, last_login, created_at, updated_at
		FROM users
		WHERE phone = $1 AND is_active = true
	`

	row := r.db.QueryRow(query, phone)
	return r.scanUser(row)
}

func (r *UserRepository) FindByPhoneOrEmail(phone, email string) (*models.User, error) {
	query := `
		SELECT id, phone, email, full_name, password_hash, is_verified, 
		       is_active, last_login, created_at, updated_at
		FROM users
		WHERE (phone = $1 OR email = $2) AND is_active = true
		LIMIT 1
	`

	row := r.db.QueryRow(query, phone, email)
	return r.scanUser(row)
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (phone, email, full_name, password_hash, is_verified, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Phone,
		user.Email,
		user.FullName,
		user.PasswordHash,
		user.IsVerified,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) UpdateLastLogin(userID int, loginTime time.Time) error {
	query := `UPDATE users SET last_login = $1, updated_at = $2 WHERE id = $3`

	_, err := r.db.Exec(query, loginTime, time.Now(), userID)
	return err
}

func (r *UserRepository) UpdateProfile(userID int, updates map[string]interface{}) error {
	// Start building the query
	query := "UPDATE users SET updated_at = $1"
	params := []interface{}{time.Now()}
	paramIndex := 2

	// Dynamically add fields to update
	if fullName, ok := updates["full_name"]; ok {
		query += fmt.Sprintf(", full_name = $%d", paramIndex)
		params = append(params, fullName)
		paramIndex++
	}
	if email, ok := updates["email"]; ok {
		query += fmt.Sprintf(", email = $%d", paramIndex)
		params = append(params, email)
		paramIndex++
	}
	if phone, ok := updates["phone"]; ok {
		query += fmt.Sprintf(", phone = $%d", paramIndex)
		params = append(params, phone)
		paramIndex++
	}

	// Add WHERE clause
	query += fmt.Sprintf(" WHERE id = $%d", paramIndex)
	params = append(params, userID)

	_, err := r.db.Exec(query, params...)
	return err
}

func (r *UserRepository) scanUser(row *sql.Row) (*models.User, error) {
	var user models.User
	var email sql.NullString
	var lastLogin sql.NullTime

	err := row.Scan(
		&user.ID,
		&user.Phone,
		&email,
		&user.FullName,
		&user.PasswordHash,
		&user.IsVerified,
		&user.IsActive,
		&lastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	if email.Valid {
		user.Email = &email.String
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return &user, nil
}

// package repository

// import (
// 	"database/sql"
// 	"fmt"
// 	"time"

// 	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
// )

// type UserRepository struct {
// 	db *sql.DB
// }

// func NewUserRepository(db *sql.DB) *UserRepository {
// 	return &UserRepository{db: db}
// }

// func (r *UserRepository) FindByPhone(phone string) (*models.User, error) {
// 	query := `
// 		SELECT id, phone, email, full_name, password_hash, is_verified,
// 		       is_active, last_login, created_at, updated_at
// 		FROM users
// 		WHERE phone = $1 AND is_active = true
// 	`

// 	row := r.db.QueryRow(query, phone)
// 	return r.scanUser(row)
// }

// func (r *UserRepository) FindByPhoneOrEmail(phone, email string) (*models.User, error) {
// 	query := `
// 		SELECT id, phone, email, full_name, password_hash, is_verified,
// 		       is_active, last_login, created_at, updated_at
// 		FROM users
// 		WHERE (phone = $1 OR email = $2) AND is_active = true
// 		LIMIT 1
// 	`

// 	row := r.db.QueryRow(query, phone, email)
// 	return r.scanUser(row)
// }

// func (r *UserRepository) Create(user *models.User) (*models.User, error) {
// 	query := `
// 		INSERT INTO users (phone, email, full_name, password_hash, is_verified, is_active, created_at, updated_at)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
// 		RETURNING id, created_at, updated_at
// 	`

// 	err := r.db.QueryRow(
// 		query,
// 		user.Phone,
// 		user.Email,
// 		user.FullName,
// 		user.PasswordHash,
// 		user.IsVerified,
// 		user.IsActive,
// 		user.CreatedAt,
// 		user.UpdatedAt,
// 	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create user: %w", err)
// 	}

// 	return user, nil
// }

// func (r *UserRepository) UpdateLastLogin(userID int, loginTime time.Time) error {
// 	query := `UPDATE users SET last_login = $1, updated_at = $2 WHERE id = $3`

// 	_, err := r.db.Exec(query, loginTime, time.Now(), userID)
// 	return err
// }

// func (r *UserRepository) scanUser(row *sql.Row) (*models.User, error) {
// 	var user models.User
// 	var email sql.NullString
// 	var lastLogin sql.NullTime

// 	err := row.Scan(
// 		&user.ID,
// 		&user.Phone,
// 		&email,
// 		&user.FullName,
// 		&user.PasswordHash,
// 		&user.IsVerified,
// 		&user.IsActive,
// 		&lastLogin,
// 		&user.CreatedAt,
// 		&user.UpdatedAt,
// 	)

// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to scan user: %w", err)
// 	}

// 	if email.Valid {
// 		user.Email = &email.String
// 	}
// 	if lastLogin.Valid {
// 		user.LastLogin = &lastLogin.Time
// 	}

// 	return &user, nil
// }
