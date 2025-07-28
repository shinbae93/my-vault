package models

import (
	"time"
)

// Secret represents a stored secret in the vault
type Secret struct {
	ID             string    `json:"id" db:"id"`
	Title          string    `json:"title" db:"title"`
	Type           string    `json:"type" db:"type"`
	EncryptedValue []byte    `json:"-" db:"encrypted_value"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// CreateSecretRequest represents the request to create a new secret
type CreateSecretRequest struct {
	Title string `json:"title" validate:"required"`
	Type  string `json:"type" validate:"required"`
	Value string `json:"value" validate:"required"`
}

// UpdateSecretRequest represents the request to update an existing secret
type UpdateSecretRequest struct {
	Title string `json:"title" validate:"required"`
	Type  string `json:"type" validate:"required"`
	Value string `json:"value" validate:"required"`
}

// SecretResponse represents the response when returning a secret
type SecretResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} 