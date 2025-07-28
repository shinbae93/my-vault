package services

import (
	"fmt"
	"sync"
	"time"

	"my-vault/internal/utils"
)

// VaultService manages the vault state and encryption key
type VaultService struct {
	mu           sync.RWMutex
	key          []byte
	salt         []byte
	isUnlocked   bool
	lastActivity time.Time
	autoLockTime time.Duration
	stopAutoLock chan struct{}
}

// NewVaultService creates a new vault service instance
func NewVaultService() *VaultService {
	return &VaultService{
		autoLockTime: 15 * time.Minute,
		stopAutoLock: make(chan struct{}),
	}
}

// Unlock unlocks the vault with the provided master password
func (v *VaultService) Unlock(masterPassword string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Generate salt if not exists (first time unlock)
	if v.salt == nil {
		salt, err := utils.GenerateSalt()
		if err != nil {
			return fmt.Errorf("failed to generate salt: %w", err)
		}
		v.salt = salt
	}

	// Derive key from master password
	key := utils.DeriveKey(masterPassword, v.salt)
	
	// Store the key in memory
	v.key = key
	v.isUnlocked = true
	v.lastActivity = time.Now()

	// Start auto-lock timer
	go v.startAutoLockTimer()

	return nil
}

// Lock locks the vault and clears the encryption key from memory
func (v *VaultService) Lock() {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Clear the key from memory
	v.key = nil
	v.isUnlocked = false

	// Stop auto-lock timer
	select {
	case v.stopAutoLock <- struct{}{}:
	default:
	}
}

// IsUnlocked returns whether the vault is currently unlocked
func (v *VaultService) IsUnlocked() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.isUnlocked
}

// GetKey returns the current encryption key
func (v *VaultService) GetKey() ([]byte, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if !v.isUnlocked {
		return nil, fmt.Errorf("vault is locked")
	}

	// Update last activity
	v.lastActivity = time.Now()

	return v.key, nil
}

// GetSalt returns the salt used for key derivation
func (v *VaultService) GetSalt() []byte {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.salt
}

// SetSalt sets the salt (used when restoring from backup)
func (v *VaultService) SetSalt(salt []byte) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.salt = salt
}

// startAutoLockTimer starts a timer that will automatically lock the vault after inactivity
func (v *VaultService) startAutoLockTimer() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			v.mu.RLock()
			if time.Since(v.lastActivity) >= v.autoLockTime {
				v.mu.RUnlock()
				v.Lock()
				return
			}
			v.mu.RUnlock()
		case <-v.stopAutoLock:
			return
		}
	}
}

// GetStatus returns the current vault status
func (v *VaultService) GetStatus() map[string]interface{} {
	v.mu.RLock()
	defer v.mu.RUnlock()

	status := map[string]interface{}{
		"unlocked": v.isUnlocked,
	}

	if v.isUnlocked {
		status["last_activity"] = v.lastActivity
		status["auto_lock_in"] = v.autoLockTime - time.Since(v.lastActivity)
	}

	return status
} 