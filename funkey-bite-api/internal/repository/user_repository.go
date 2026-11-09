package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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

// SetPassword gives an account a password; guest customers created at checkout
// have none until they register.
// phoneDigitsKey is the part of a phone number customers are matched on:
// its last ten digits, so 0803 111 2222 and +234 803 111 2222 match.
func phoneDigitsKey(phone string) string {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)
	if len(digits) > 10 {
		return digits[len(digits)-10:]
	}
	return digits
}

// phoneMatch matches users.phone against $1 (phoneDigitsKey) or, for a short
// number, exactly against $2.
const phoneMatch = `((length($1) = 10 AND right(regexp_replace(phone, '\D', '', 'g'), 10) = $1) OR phone = $2)`

// FindByPhoneAnyFormat returns the customer with this number in any format,
// active or not, or nil.
func (r *UserRepository) FindByPhoneAnyFormat(phone string) (*models.User, error) {
	row := r.db.QueryRow(`
		SELECT id, phone, email, full_name, password_hash, is_verified,
		       is_active, last_login, created_at, updated_at
		FROM users WHERE `+phoneMatch+` ORDER BY id LIMIT 1
	`, phoneDigitsKey(phone), phone)
	return r.scanUser(row)
}

// FindByEmail returns the customer with this email, ignoring case, or nil.
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	row := r.db.QueryRow(`
		SELECT id, phone, email, full_name, password_hash, is_verified,
		       is_active, last_login, created_at, updated_at
		FROM users WHERE lower(email) = lower($1) ORDER BY id LIMIT 1
	`, email)
	return r.scanUser(row)
}

// LatestOrderEmail returns the email on a customer's most recent order, for
// a guest record saved without one.
func (r *UserRepository) LatestOrderEmail(userID int) (string, error) {
	var email string
	err := r.db.QueryRow(`
		SELECT customer_email FROM orders
		WHERE user_id = $1 AND customer_email IS NOT NULL AND customer_email <> ''
		ORDER BY created_at DESC LIMIT 1
	`, userID).Scan(&email)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return email, err
}

// ClaimGuest turns a guest record into an account once its owner has proved
// the email, setting their password, name and email. It returns false if the
// record already has a password (claimed in the meantime).
func (r *UserRepository) ClaimGuest(userID int, fullName, email, passwordHash string) (bool, error) {
	result, err := r.db.Exec(`
		UPDATE users
		SET password_hash = $2, full_name = $3, email = $4, is_verified = true, updated_at = NOW()
		WHERE id = $1 AND password_hash = ''
	`, userID, passwordHash, fullName, email)
	if err != nil {
		return false, fmt.Errorf("failed to claim customer record: %w", err)
	}
	n, err := result.RowsAffected()
	return n > 0, err
}

// FindOrCreateGuestTx returns the customer an order should be linked to,
// creating a passwordless record inside tx for a first-time guest, so the
// admin sees every customer and they can register later. Numbers match in any
// format. A number that belongs to a registered account is only linked when
// the checkout email is that account's email; otherwise it returns 0 and the
// order stays unlinked, so nobody can collect another person's orders by
// registering their number. If the email already belongs to another customer,
// the record is created without it. Conflicts are skipped rather than raised,
// since a failed statement would abort the caller's transaction.
func (r *UserRepository) FindOrCreateGuestTx(ctx context.Context, tx *sql.Tx, name, phone string, email *string) (int, error) {
	find := func() (int, bool, error) {
		var id int
		var passwordHash string
		var existingEmail sql.NullString
		err := tx.QueryRowContext(ctx, `
			SELECT id, password_hash, email FROM users WHERE `+phoneMatch+` ORDER BY id LIMIT 1
		`, phoneDigitsKey(phone), phone).Scan(&id, &passwordHash, &existingEmail)
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		if err != nil {
			return 0, false, fmt.Errorf("failed to find customer: %w", err)
		}
		if passwordHash == "" {
			return id, true, nil
		}
		if email != nil && existingEmail.Valid && strings.EqualFold(existingEmail.String, *email) {
			return id, true, nil
		}
		return 0, true, nil
	}

	if id, found, err := find(); err != nil || found {
		return id, err
	}

	insert := `
		INSERT INTO users (phone, email, full_name, password_hash, is_verified, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, '', false, true, NOW(), NOW())
		ON CONFLICT DO NOTHING
		RETURNING id
	`
	for _, candidate := range []*string{email, nil} {
		var id int
		err := tx.QueryRowContext(ctx, insert, phone, candidate, name).Scan(&id)
		if err == nil {
			return id, nil
		}
		if err != sql.ErrNoRows {
			return 0, fmt.Errorf("failed to create customer: %w", err)
		}
		// Nothing inserted: the number was taken meanwhile (use that
		// customer) or only the email is (try again without it).
		if id, found, err := find(); err != nil || found {
			return id, err
		}
		if candidate == nil {
			break
		}
	}
	return 0, fmt.Errorf("failed to create customer for phone ending %s", lastFour(phone))
}

func lastFour(phone string) string {
	if len(phone) <= 4 {
		return phone
	}
	return phone[len(phone)-4:]
}

func (r *UserRepository) SetPassword(userID int, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`

	_, err := r.db.Exec(query, passwordHash, time.Now(), userID)
	return err
}

func (r *UserRepository) UpdateLastLogin(userID int, loginTime time.Time) error {
	query := `UPDATE users SET last_login = $1, updated_at = $2 WHERE id = $3`

	_, err := r.db.Exec(query, loginTime, time.Now(), userID)
	return err
}

func (r *UserRepository) UpdateProfile(userID int, updates map[string]interface{}) error {
	query := "UPDATE users SET updated_at = $1"
	params := []interface{}{time.Now()}
	paramIndex := 2

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
