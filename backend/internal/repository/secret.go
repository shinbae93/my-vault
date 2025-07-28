package repository

import (
	"context"
	"fmt"
	"time"

	"my-vault/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SecretRepository handles database operations for secrets
type SecretRepository struct {
	pool *pgxpool.Pool
}

// NewSecretRepository creates a new secret repository
func NewSecretRepository(db *PostgresDB) *SecretRepository {
	return &SecretRepository{
		pool: db.GetPool(),
	}
}

// Create creates a new secret in the database
func (r *SecretRepository) Create(ctx context.Context, secret *models.Secret) error {
	query := `
		INSERT INTO secrets (id, title, type, encrypted_value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	secret.ID = uuid.New().String()
	now := time.Now()
	secret.CreatedAt = now
	secret.UpdatedAt = now

	_, err := r.pool.Exec(ctx, query,
		secret.ID,
		secret.Title,
		secret.Type,
		secret.EncryptedValue,
		secret.CreatedAt,
		secret.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	return nil
}

// Get retrieves a secret by ID
func (r *SecretRepository) Get(ctx context.Context, id string) (*models.Secret, error) {
	query := `
		SELECT id, title, type, encrypted_value, created_at, updated_at
		FROM secrets
		WHERE id = $1
	`

	var secret models.Secret
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&secret.ID,
		&secret.Title,
		&secret.Type,
		&secret.EncryptedValue,
		&secret.CreatedAt,
		&secret.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	return &secret, nil
}

// List retrieves all secrets
func (r *SecretRepository) List(ctx context.Context) ([]*models.Secret, error) {
	query := `
		SELECT id, title, type, encrypted_value, created_at, updated_at
		FROM secrets
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	defer rows.Close()

	var secrets []*models.Secret
	for rows.Next() {
		var secret models.Secret
		err := rows.Scan(
			&secret.ID,
			&secret.Title,
			&secret.Type,
			&secret.EncryptedValue,
			&secret.CreatedAt,
			&secret.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan secret: %w", err)
		}
		secrets = append(secrets, &secret)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating secrets: %w", err)
	}

	return secrets, nil
}

// Update updates an existing secret
func (r *SecretRepository) Update(ctx context.Context, secret *models.Secret) error {
	query := `
		UPDATE secrets
		SET title = $1, type = $2, encrypted_value = $3, updated_at = $4
		WHERE id = $5
	`

	secret.UpdatedAt = time.Now()

	result, err := r.pool.Exec(ctx, query,
		secret.Title,
		secret.Type,
		secret.EncryptedValue,
		secret.UpdatedAt,
		secret.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("secret not found")
	}

	return nil
}

// Delete removes a secret by ID
func (r *SecretRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM secrets WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("secret not found")
	}

	return nil
} 