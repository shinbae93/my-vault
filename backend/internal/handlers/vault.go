package handlers

import (
	"net/http"

	"my-vault/internal/models"
	"my-vault/internal/services"

	"github.com/gin-gonic/gin"
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

// Unlock unlocks the vault with the provided master password
// @Summary Unlock vault
// @Description Unlock the vault using the master password
// @Tags vault
// @Accept json
// @Produce json
// @Param request body models.UnlockRequest true "Unlock request"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/unlock [post]
func (h *VaultHandler) Unlock(c *gin.Context) {
	var req models.UnlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: "Failed to parse request body",
		})
		return
	}

	if req.MasterPassword == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: "Master password is required",
		})
		return
	}

	if err := h.vaultService.Unlock(req.MasterPassword); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Authentication failed",
			Message: "Failed to unlock vault",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Vault unlocked successfully",
	})
}

// Lock locks the vault
// @Summary Lock vault
// @Description Lock the vault and clear encryption key from memory
// @Tags vault
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Router /api/lock [post]
func (h *VaultHandler) Lock(c *gin.Context) {
	h.vaultService.Lock()

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Vault locked successfully",
	})
}

// Status returns the current vault status
// @Summary Get vault status
// @Description Get the current status of the vault
// @Tags vault
// @Produce json
// @Success 200 {object} models.VaultStatus
// @Router /api/status [get]
func (h *VaultHandler) Status(c *gin.Context) {
	status := h.vaultService.GetStatus()

	c.JSON(http.StatusOK, status)
}

// RequireUnlocked is middleware that ensures the vault is unlocked
func (h *VaultHandler) RequireUnlocked() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !h.vaultService.IsUnlocked() {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "Vault is locked",
				Message: "The vault must be unlocked before accessing secrets",
			})
			c.Abort()
			return
		}
		c.Next()
	}
} 