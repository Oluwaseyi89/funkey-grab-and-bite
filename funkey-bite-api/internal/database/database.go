package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"funkey-grab-and-bite/funkey-bite-api/internal/utils"

	_ "github.com/lib/pq"
)

func InitializeDatabase() *sql.DB {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "funkey_grab_bite")
	sslmode := getEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * 60)

	log.Println("✅ Database connection established successfully")

	if err := runMigrations(db); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	return db
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			phone VARCHAR(20) UNIQUE NOT NULL,
			email VARCHAR(200) UNIQUE,
			full_name VARCHAR(200) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			is_verified BOOLEAN DEFAULT true,
			is_active BOOLEAN DEFAULT true,
			last_login TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS menu_categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			display_order INTEGER,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS menu_items (
			id SERIAL PRIMARY KEY,
			category_id INTEGER REFERENCES menu_categories(id),
			name VARCHAR(200) NOT NULL,
			description TEXT,
			price DECIMAL(10,2) NOT NULL,
			image_url VARCHAR(500),
			is_available BOOLEAN DEFAULT true,
			is_pre_order BOOLEAN DEFAULT false,
			preparation_time INTEGER,
			tags JSONB DEFAULT '[]'::jsonb,
			nutritional_info JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			order_number VARCHAR(50) UNIQUE NOT NULL,
			user_id INTEGER REFERENCES users(id),
			customer_name VARCHAR(200) NOT NULL,
			customer_phone VARCHAR(20) NOT NULL,
			customer_email VARCHAR(200),
			order_type VARCHAR(20) NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			total_amount DECIMAL(10,2) NOT NULL,
			notes TEXT,
			pickup_time TIMESTAMP,
			estimated_ready_time TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
			menu_item_id INTEGER NOT NULL,
			name VARCHAR(200) NOT NULL,
			quantity INTEGER NOT NULL,
			unit_price DECIMAL(10,2) NOT NULL,
			special_instructions TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS catering_requests (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			event_name VARCHAR(200),
			contact_name VARCHAR(200) NOT NULL,
			contact_phone VARCHAR(20) NOT NULL,
			contact_email VARCHAR(200),
			event_date DATE NOT NULL,
			event_time TIME,
			guest_count INTEGER,
			event_type VARCHAR(100),
			budget DECIMAL(10,2),
			special_requests TEXT,
			status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS admin_users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(100) UNIQUE NOT NULL,
			email VARCHAR(200) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(50) DEFAULT 'manager',
			is_active BOOLEAN DEFAULT true,
			last_login TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS business_settings (
			id SERIAL PRIMARY KEY,
			business_name VARCHAR(200) NOT NULL,
			phone_number VARCHAR(20) NOT NULL,
			email VARCHAR(200) NOT NULL,
			address TEXT NOT NULL,
			opening_hours TEXT NOT NULL,
			delivery_fee DECIMAL(10,2) DEFAULT 2.99,
			min_order_amount DECIMAL(10,2) DEFAULT 10.00,
			tax_rate DECIMAL(5,2) DEFAULT 8.5,
			is_delivery_open BOOLEAN DEFAULT true,
			is_pickup_open BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS promotions (
			id SERIAL PRIMARY KEY,
			code VARCHAR(50) UNIQUE NOT NULL,
			title VARCHAR(200) NOT NULL,
			description TEXT,
			promotion_type VARCHAR(20) NOT NULL,
			discount_value DECIMAL(10,2) NOT NULL,
			max_discount DECIMAL(10,2),
			min_order_amount DECIMAL(10,2),
			valid_from TIMESTAMP NOT NULL,
			valid_until TIMESTAMP NOT NULL,
			usage_limit INTEGER,
			used_count INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS promotion_usage (
			id SERIAL PRIMARY KEY,
			promotion_id INTEGER REFERENCES promotions(id),
			order_id INTEGER REFERENCES orders(id),
			customer_id INTEGER REFERENCES users(id),
			discount_applied DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS inventory_items (
			id SERIAL PRIMARY KEY,
			menu_item_id INTEGER REFERENCES menu_items(id) UNIQUE,
			name VARCHAR(200) NOT NULL,
			current_stock INTEGER NOT NULL DEFAULT 0,
			minimum_stock INTEGER DEFAULT 10,
			reorder_point INTEGER DEFAULT 5,
			unit VARCHAR(20) DEFAULT 'pieces',
			is_active BOOLEAN DEFAULT true,
			last_restocked TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS inventory_history (
			id SERIAL PRIMARY KEY,
			inventory_item_id INTEGER REFERENCES inventory_items(id),
			previous_stock INTEGER NOT NULL,
			new_stock INTEGER NOT NULL,
			change INTEGER NOT NULL,
			operation VARCHAR(20) NOT NULL,
			reason VARCHAR(200) NOT NULL,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS inventory_alerts (
			id SERIAL PRIMARY KEY,
			inventory_item_id INTEGER REFERENCES inventory_items(id),
			alert_type VARCHAR(20) NOT NULL,
			message TEXT NOT NULL,
			is_resolved BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			resolved_at TIMESTAMP,
			UNIQUE(inventory_item_id, alert_type)
		)`,

		`CREATE TABLE IF NOT EXISTS notifications (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			type VARCHAR(50) NOT NULL,
			title VARCHAR(200) NOT NULL,
			message TEXT NOT NULL,
			is_read BOOLEAN DEFAULT false,
			reference_id INTEGER,
			reference_type VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			read_at TIMESTAMP
		)`,

		`ALTER TABLE menu_items 
		ADD COLUMN IF NOT EXISTS search_vector tsvector GENERATED ALWAYS AS (
			setweight(to_tsvector('english', COALESCE(name, '')), 'A') ||
			setweight(to_tsvector('english', COALESCE(description, '')), 'B') ||
			setweight(to_tsvector('english', COALESCE(tags::text, '')), 'C')
		) STORED`,

		`CREATE INDEX IF NOT EXISTS idx_menu_items_search 
		ON menu_items USING GIN(search_vector)`,
	}

	if err := runMigrationsWithStatements(db, migrations); err != nil {
		return err
	}

	if err := ensureDefaultAdminUser(db); err != nil {
		return fmt.Errorf("default admin bootstrap failed: %w", err)
	}

	log.Println("✅ Database migrations completed")
	return nil
}

func runMigrationsWithStatements(db *sql.DB, migrations []string) error {
	for i, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	return nil
}

func ensureDefaultAdminUser(db *sql.DB) error {
	defaultEmail := getEnv("DEFAULT_ADMIN_EMAIL", "admin@funkey.com")
	defaultUsername := getEnv("DEFAULT_ADMIN_USERNAME", "admin")
	defaultPassword := getEnv("DEFAULT_ADMIN_PASSWORD", "admin123")
	defaultRole := getEnv("DEFAULT_ADMIN_ROLE", "admin")

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM admin_users WHERE email = $1)", defaultEmail).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed checking default admin existence: %w", err)
	}

	if exists {
		return nil
	}

	hashedPassword, err := utils.HashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("failed hashing default admin password: %w", err)
	}

	_, err = db.Exec(
		`INSERT INTO admin_users (username, email, password_hash, role, is_active) VALUES ($1, $2, $3, $4, true)`,
		defaultUsername,
		defaultEmail,
		hashedPassword,
		defaultRole,
	)
	if err != nil {
		return fmt.Errorf("failed creating default admin user: %w", err)
	}

	log.Printf("✅ Created default admin user: %s", defaultEmail)
	return nil
}

func CloseDatabase(db *sql.DB) {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}
