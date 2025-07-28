package services

import (
	"context"
	"fmt"

	"my-vault/internal/models"
	"my-vault/internal/repository"
	"my-vault/internal/utils"
)

// SecretService handles business logic for secrets
type SecretService struct {
	repo         *repository.SecretRepository
	vaultService *VaultService
}

// NewSecretService creates a new secret service
func NewSecretService(repo *repository.SecretRepository, vaultService *VaultService) *SecretService {
	return &SecretService{
		repo:         repo,
		vaultService: vaultService,
	}
}

// Create creates a new secret
func (s *SecretService) Create(ctx context.Context, req *models.CreateSecretRequest) (*models.SecretResponse, error) {
	// Get encryption key from vault
	key, err := s.vaultService.GetKey()
	if err != nil {
		return nil, fmt.Errorf("vault is locked: %w", err)
	}

	// Encrypt the secret value
	encryptedValue, err := utils.Encrypt([]byte(req.Value), key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// Create secret model
	secret := &models.Secret{
		Title:          req.Title,
		Type:           req.Type,
		EncryptedValue: encryptedValue,
	}

	// Save to database
	if err := s.repo.Create(ctx, secret); err != nil {
		return nil, fmt.Errorf("failed to save secret: %w", err)
	}

	// Return response
	return &models.SecretResponse{
		ID:        secret.ID,
		Title:     secret.Title,
		Type:      secret.Type,
		Value:     req.Value, // Return decrypted value
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	}, nil
}

// Get retrieves a secret by ID
func (s *SecretService) Get(ctx context.Context, id string) (*models.SecretResponse, error) {
	// Get encryption key from vault
	key, err := s.vaultService.GetKey()
	if err != nil {
		return nil, fmt.Errorf("vault is locked: %w", err)
	}

	// Get secret from database
	secret, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	// Decrypt the secret value
	decryptedValue, err := utils.Decrypt(secret.EncryptedValue, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// Return response
	return &models.SecretResponse{
		ID:        secret.ID,
		Title:     secret.Title,
		Type:      secret.Type,
		Value:     string(decryptedValue),
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	}, nil
}

// List retrieves all secrets
func (s *SecretService) List(ctx context.Context) ([]*models.SecretResponse, error) {
	// Get encryption key from vault
	key, err := s.vaultService.GetKey()
	if err != nil {
		return nil, fmt.Errorf("vault is locked: %w", err)
	}

	// Get secrets from database
	secrets, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	// Decrypt and convert to responses
	var responses []*models.SecretResponse
	for _, secret := range secrets {
		decryptedValue, err := utils.Decrypt(secret.EncryptedValue, key)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt secret %s: %w", secret.ID, err)
		}

		responses = append(responses, &models.SecretResponse{
			ID:        secret.ID,
			Title:     secret.Title,
			Type:      secret.Type,
			Value:     string(decryptedValue),
			CreatedAt: secret.CreatedAt,
			UpdatedAt: secret.UpdatedAt,
		})
	}

	return responses, nil
}

// Update updates an existing secret
func (s *SecretService) Update(ctx context.Context, id string, req *models.UpdateSecretRequest) (*models.SecretResponse, error) {
	// Get encryption key from vault
	key, err := s.vaultService.GetKey()
	if err != nil {
		return nil, fmt.Errorf("vault is locked: %w", err)
	}

	// Get existing secret
	secret, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	// Encrypt the new secret value
	encryptedValue, err := utils.Encrypt([]byte(req.Value), key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// Update secret fields
	secret.Title = req.Title
	secret.Type = req.Type
	secret.EncryptedValue = encryptedValue

	// Save to database
	if err := s.repo.Update(ctx, secret); err != nil {
		return nil, fmt.Errorf("failed to update secret: %w", err)
	}

	// Return response
	return &models.SecretResponse{
		ID:        secret.ID,
		Title:     secret.Title,
		Type:      secret.Type,
		Value:     req.Value, // Return decrypted value
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	}, nil
}

// Delete removes a secret
func (s *SecretService) Delete(ctx context.Context, id string) error {
	// Check if vault is unlocked (we don't need the key for deletion)
	if !s.vaultService.IsUnlocked() {
		return fmt.Errorf("vault is locked")
	}

	return s.repo.Delete(ctx, id)
} 