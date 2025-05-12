package webauthentication

import (
	"msnserver/internal/web/auth"
	"msnserver/pkg/database"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PasskeyResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (wc *WebauthnController) GetPasskeys(c echo.Context) error {
	// Get user email from JWT
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*auth.JwtCustomClaims)
	email := claims.Subject

	// Retrieve user from database
	var user database.User
	if err := wc.db.Preload("WebauthnCredentials").Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	// Prepare response
	passkeyResponses := make([]PasskeyResponse, len(user.WebauthnCredentials))
	for i, passkey := range user.WebauthnCredentials {
		passkeyResponses[i] = PasskeyResponse{
			ID:   passkey.ID,
			Name: passkey.Name,
		}
	}

	return c.JSON(http.StatusOK, passkeyResponses)
}

func (wc *WebauthnController) DeletePasskey(c echo.Context) error {
	// Get user email from JWT
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*auth.JwtCustomClaims)
	email := claims.Subject

	// Retrieve user from database
	var user database.User
	if err := wc.db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	// Get passkey ID from URL parameter
	passkeyID := c.Param("id")

	// Delete the passkey from the database
	query := wc.db.Where("id = ? AND user_id = ?", passkeyID, user.ID).Delete(&database.Credential{})
	if query.Error != nil {
		return query.Error
	}
	if query.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "passkey not found")
	}

	return c.NoContent(http.StatusNoContent)
}
