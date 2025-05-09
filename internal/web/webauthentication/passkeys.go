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

func (uc *WebauthnController) GetPasskeys(c echo.Context) error {
	// Get user email from JWT
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*auth.JwtCustomClaims)
	email := claims.Subject

	// Retrieve user from database
	var user database.User
	if err := uc.db.Preload("WebauthnCredentials").Where("email = ?", email).First(&user).Error; err != nil {
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
