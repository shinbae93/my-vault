package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresDB wraps the database connection pool
type PostgresDB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB() (*PostgresDB, error) {
	// Get database configuration from environment
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "vaultbox")
	password := getEnv("DB_PASSWORD", "supersecret")
	dbname := getEnv("DB_NAME", "vaultbox")

	// Build connection string
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	// Create connection pool
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Initialize database schema
	if err := initSchema(pool); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &PostgresDB{pool: pool}, nil
}

// Close closes the database connection pool
func (db *PostgresDB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

// GetPool returns the underlying connection pool
func (db *PostgresDB) GetPool() *pgxpool.Pool {
	return db.pool
}

// initSchema creates the necessary database tables
func initSchema(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create secrets table
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS secrets (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(255) NOT NULL,
			type VARCHAR(100) NOT NULL,
			encrypted_value BYTEA NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
		
		-- Create index on title for faster searches
		CREATE INDEX IF NOT EXISTS idx_secrets_title ON secrets(title);
		
		-- Create index on type for filtering
		CREATE INDEX IF NOT EXISTS idx_secrets_type ON secrets(type);
	`

	_, err := pool.Exec(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create secrets table: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 