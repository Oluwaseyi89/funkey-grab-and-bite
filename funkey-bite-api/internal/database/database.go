package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// InitializeDatabase creates and returns a database connection
func InitializeDatabase() *sql.DB {
	// Get database connection string from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "funkey_grab_bite")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * 60) // 5 minutes

	log.Println("✅ Database connection established successfully")

	// Run migrations (optional)
	runMigrations(db)

	return db
}

// getEnv helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// runMigrations creates database tables if they don't exist
func runMigrations(db *sql.DB) {
	// Create tables if they don't exist
	migrations := []string{
		// Users table
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

		// Menu categories
		`CREATE TABLE IF NOT EXISTS menu_categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			display_order INTEGER,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Menu items
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

		// Orders table
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

		// Order items
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

		// Catering requests
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

		// Admin users (for admin dashboard)
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
	}

	// Execute migrations
	for i, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			log.Printf("Warning: Migration %d failed: %v", i+1, err)
		}
	}

	log.Println("✅ Database migrations completed")
}

// CloseDatabase safely closes the database connection
func CloseDatabase(db *sql.DB) {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}
