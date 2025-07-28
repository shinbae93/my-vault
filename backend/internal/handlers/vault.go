package handlers

import (
	"encoding/json"
	"net/http"

	"my-vault/internal/services"
)

// VaultHandler handles vault-related HTTP requests
type VaultHandler struct {
	vaultService *services.VaultService
}

// NewVaultHandler creates a new vault handler
func NewVaultHandler(vaultService *services.VaultService) *VaultHandler {
	return &VaultHandler{
		vaultService: vaultService,
	}
}

// UnlockRequest represents the request to unlock the vault
type UnlockRequest struct {
	MasterPassword string `json:"master_password" validate:"required"`
}

// Unlock unlocks the vault with the provided master password
func (h *VaultHandler) Unlock(w http.ResponseWriter, r *http.Request) {
	var req UnlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.MasterPassword == "" {
		http.Error(w, "Master password is required", http.StatusBadRequest)
		return
	}

	if err := h.vaultService.Unlock(req.MasterPassword); err != nil {
		http.Error(w, "Failed to unlock vault", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vault unlocked successfully",
	})
}

// Lock locks the vault
func (h *VaultHandler) Lock(w http.ResponseWriter, r *http.Request) {
	h.vaultService.Lock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vault locked successfully",
	})
}

// Status returns the current vault status
func (h *VaultHandler) Status(w http.ResponseWriter, r *http.Request) {
	status := h.vaultService.GetStatus()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// RequireUnlocked is middleware that ensures the vault is unlocked
func (h *VaultHandler) RequireUnlocked(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.vaultService.IsUnlocked() {
			http.Error(w, "Vault is locked", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
} 