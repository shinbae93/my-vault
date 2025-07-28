package handlers

import (
	"net/http"

	"my-vault/internal/models"
	"my-vault/internal/services"

	"github.com/gin-gonic/gin"
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
// @Summary List all secrets
// @Description Retrieve all secrets from the vault
// @Tags secrets
// @Produce json
// @Success 200 {array} models.SecretResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/secrets [get]
func (h *SecretHandler) List(c *gin.Context) {
	secrets, err := h.secretService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to list secrets",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, secrets)
}

// Create creates a new secret
// @Summary Create a new secret
// @Description Create a new secret in the vault
// @Tags secrets
// @Accept json
// @Produce json
// @Param request body models.CreateSecretRequest true "Secret creation request"
// @Success 201 {object} models.SecretResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/secrets [post]
func (h *SecretHandler) Create(c *gin.Context) {
	var req models.CreateSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: "Failed to parse request body",
		})
		return
	}

	// Basic validation
	if req.Title == "" || req.Type == "" || req.Value == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: "Title, type, and value are required",
		})
		return
	}

	secret, err := h.secretService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create secret",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, secret)
}

// Get retrieves a secret by ID
// @Summary Get a secret by ID
// @Description Retrieve a specific secret by its ID
// @Tags secrets
// @Produce json
// @Param id path string true "Secret ID"
// @Success 200 {object} models.SecretResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/secrets/{id} [get]
func (h *SecretHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Secret ID is required",
		})
		return
	}

	secret, err := h.secretService.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Secret not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, secret)
}

// Update updates an existing secret
// @Summary Update a secret
// @Description Update an existing secret by its ID
// @Tags secrets
// @Accept json
// @Produce json
// @Param id path string true "Secret ID"
// @Param request body models.UpdateSecretRequest true "Secret update request"
// @Success 200 {object} models.SecretResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/secrets/{id} [put]
func (h *SecretHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Secret ID is required",
		})
		return
	}

	var req models.UpdateSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: "Failed to parse request body",
		})
		return
	}

	// Basic validation
	if req.Title == "" || req.Type == "" || req.Value == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: "Title, type, and value are required",
		})
		return
	}

	secret, err := h.secretService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update secret",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, secret)
}

// Delete removes a secret
// @Summary Delete a secret
// @Description Delete a secret by its ID
// @Tags secrets
// @Param id path string true "Secret ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/secrets/{id} [delete]
func (h *SecretHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Secret ID is required",
		})
		return
	}

	if err := h.secretService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete secret",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
} 