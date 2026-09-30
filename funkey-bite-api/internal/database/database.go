package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
	"funkey-grab-and-bite/funkey-bite-api/migrations"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
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

	// -----------------------------------------------------------------
	// SERVERLESS ADAPTIVE CONNECTION POOLING
	// -----------------------------------------------------------------
	// AWS Lambda automatically populates AWS_LAMBDA_FUNCTION_NAME.
	// If detected, we aggressively lower defaults to safely scale on RDS Free Tier.
	isLambda := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""

	defaultMaxOpen := "25"
	defaultMaxIdle := "5"
	if isLambda {
		defaultMaxOpen = "3"
		defaultMaxIdle = "1"
	}

	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", defaultMaxOpen))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", defaultMaxIdle))

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("✅ Database pooling configured: MaxOpen=%d, MaxIdle=%d (Lambda Detected: %t)",
		maxOpenConns, maxIdleConns, isLambda)
	// -----------------------------------------------------------------

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

// runMigrations brings the schema up to date from the versioned SQL files in
// /migrations, which are the single source of truth for the schema. Applied
// versions are tracked in schema_migrations, and golang-migrate holds a
// Postgres advisory lock while running so concurrent instances can't race.
func runMigrations(db *sql.DB) error {
	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed loading embedded migrations: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed creating migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed initialising migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed applying migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("failed reading migration version: %w", err)
	}
	log.Printf("✅ Database migrations completed (version=%d, dirty=%t)", version, dirty)

	// The default admin needs env-supplied credentials and a bcrypt hash, so it
	// is bootstrapped in code once the schema is in place rather than in SQL.
	if err := ensureDefaultAdminUser(db); err != nil {
		return fmt.Errorf("default admin bootstrap failed: %w", err)
	}

	return nil
}

func ensureDefaultAdminUser(db *sql.DB) error {
	var adminCount int
	err := db.QueryRow("SELECT COUNT(1) FROM admin_users").Scan(&adminCount)
	if err != nil {
		return fmt.Errorf("failed checking admin user count: %w", err)
	}

	// Idempotent bootstrap: never create additional default admins when one already exists.
	if adminCount > 0 {
		return nil
	}

	defaultEmail := strings.TrimSpace(os.Getenv("DEFAULT_ADMIN_EMAIL"))
	defaultUsername := strings.TrimSpace(os.Getenv("DEFAULT_ADMIN_USERNAME"))
	defaultPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	defaultRole := strings.TrimSpace(getEnv("DEFAULT_ADMIN_ROLE", "admin"))

	missingVars := make([]string, 0, 3)
	if defaultEmail == "" {
		missingVars = append(missingVars, "DEFAULT_ADMIN_EMAIL")
	}
	if defaultUsername == "" {
		missingVars = append(missingVars, "DEFAULT_ADMIN_USERNAME")
	}
	if defaultPassword == "" {
		missingVars = append(missingVars, "DEFAULT_ADMIN_PASSWORD")
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("missing required default admin bootstrap env vars: %s", strings.Join(missingVars, ", "))
	}

	environment := strings.ToLower(strings.TrimSpace(getEnv("ENVIRONMENT", "development")))
	if isProductionLikeEnvironment(environment) && usesWeakDefaultAdminCredentials(defaultEmail, defaultUsername, defaultPassword) {
		return fmt.Errorf("weak default admin credentials are not allowed when ENVIRONMENT=%s", environment)
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

func isProductionLikeEnvironment(env string) bool {
	switch env {
	case "production", "prod", "staging":
		return true
	default:
		return false
	}
}

func usesWeakDefaultAdminCredentials(email, username, password string) bool {
	weakEmails := map[string]struct{}{
		"admin@funkey.com":  {},
		"admin@example.com": {},
		"admin@localhost":   {},
	}
	weakUsernames := map[string]struct{}{
		"admin":         {},
		"administrator": {},
		"root":          {},
	}
	weakPasswords := map[string]struct{}{
		"admin":       {},
		"admin123":    {},
		"password":    {},
		"password123": {},
		"changeme":    {},
		"123456":      {},
	}

	if _, ok := weakEmails[strings.ToLower(strings.TrimSpace(email))]; ok {
		return true
	}
	if _, ok := weakUsernames[strings.ToLower(strings.TrimSpace(username))]; ok {
		return true
	}
	normalizedPassword := strings.ToLower(strings.TrimSpace(password))
	if _, ok := weakPasswords[normalizedPassword]; ok {
		return true
	}

	return len(strings.TrimSpace(password)) < 12
}

func CloseDatabase(db *sql.DB) {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}
