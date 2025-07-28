package handlers

import (
	"encoding/json"
	"net/http"

	"my-vault/internal/models"
	"my-vault/internal/services"
)

// SecretHandler handles secret-related HTTP requests
type SecretHandler struct {
	secretService *services.SecretService
	vaultService  *services.VaultService
}

// NewSecretHandler creates a new secret handler
func NewSecretHandler(secretService *services.SecretService, vaultService *services.VaultService) *SecretHandler {
	return &SecretHandler{
		secretService: secretService,
		vaultService:  vaultService,
	}
}

// List retrieves all secrets
func (h *SecretHandler) List(w http.ResponseWriter, r *http.Request) {
	secrets, err := h.secretService.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secrets)
}

// Create creates a new secret
func (h *SecretHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Title == "" || req.Type == "" || req.Value == "" {
		http.Error(w, "Title, type, and value are required", http.StatusBadRequest)
		return
	}

	secret, err := h.secretService.Create(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(secret)
}

// Get retrieves a secret by ID
func (h *SecretHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path - simplified for now
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Secret ID is required", http.StatusBadRequest)
		return
	}

	secret, err := h.secretService.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secret)
}

// Update updates an existing secret
func (h *SecretHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path - simplified for now
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Secret ID is required", http.StatusBadRequest)
		return
	}

	var req models.UpdateSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Title == "" || req.Type == "" || req.Value == "" {
		http.Error(w, "Title, type, and value are required", http.StatusBadRequest)
		return
	}

	secret, err := h.secretService.Update(r.Context(), id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secret)
}

// Delete removes a secret
func (h *SecretHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path - simplified for now
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Secret ID is required", http.StatusBadRequest)
		return
	}

	if err := h.secretService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
} 