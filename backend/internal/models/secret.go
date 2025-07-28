package models

import (
	"time"
)

// Secret represents a stored secret in the vault
// @Description Secret entity with encrypted data
type Secret struct {
	ID             string    `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title          string    `json:"title" db:"title" example:"GitHub API Token"`
	Type           string    `json:"type" db:"type" example:"api_token"`
	EncryptedValue []byte    `json:"-" db:"encrypted_value"`
	CreatedAt      time.Time `json:"created_at" db:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// CreateSecretRequest represents the request to create a new secret
// @Description Request payload for creating a new secret
type CreateSecretRequest struct {
	Title string `json:"title" validate:"required" example:"GitHub API Token" binding:"required"`
	Type  string `json:"type" validate:"required" example:"api_token" binding:"required"`
	Value string `json:"value" validate:"required" example:"ghp_xxxxxxxxxxxxxxxxxxxx" binding:"required"`
}

// UpdateSecretRequest represents the request to update an existing secret
// @Description Request payload for updating an existing secret
type UpdateSecretRequest struct {
	Title string `json:"title" validate:"required" example:"Updated GitHub Token" binding:"required"`
	Type  string `json:"type" validate:"required" example:"api_token" binding:"required"`
	Value string `json:"value" validate:"required" example:"ghp_yyyyyyyyyyyyyyyyyyyy" binding:"required"`
}

// SecretResponse represents the response when returning a secret
// @Description Response payload for secret data
type SecretResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title     string    `json:"title" example:"GitHub API Token"`
	Type      string    `json:"type" example:"api_token"`
	Value     string    `json:"value" example:"ghp_xxxxxxxxxxxxxxxxxxxx"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// UnlockRequest represents the request to unlock the vault
// @Description Request payload for unlocking the vault
type UnlockRequest struct {
	MasterPassword string `json:"master_password" validate:"required" example:"my-secure-password" binding:"required"`
}

// VaultStatus represents the current vault status
// @Description Response payload for vault status
type VaultStatus struct {
	Unlocked    bool      `json:"unlocked" example:"true"`
	LastActivity *time.Time `json:"last_activity,omitempty" example:"2024-01-15T10:30:00Z"`
	AutoLockIn  *string   `json:"auto_lock_in,omitempty" example:"14m30s"`
}

// ErrorResponse represents an error response
// @Description Error response payload
type ErrorResponse struct {
	Error   string `json:"error" example:"Vault is locked"`
	Message string `json:"message" example:"The vault must be unlocked before accessing secrets"`
}

// SuccessResponse represents a success response
// @Description Success response payload
type SuccessResponse struct {
	Message string `json:"message" example:"Vault unlocked successfully"`
} 