package user

import (
	"msnserver/internal/web/auth"
	"msnserver/pkg/database"
	"msnserver/pkg/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PasswordBody struct {
	OldPassword string `json:"old_password" validate:"required,min=8,max=16,alphanum"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=16,alphanum"`
}

// UpdatePassword godoc
//
//	@Summary		Update password
//	@Description	Update user password
//	@Tags			user
//	@Accept			json
//	@Produce		plain
//	@Param			password	body		PasswordBody	true	"Password information"
//	@Success		200			{object}	UserResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		500			{string}	string	"internal server error"
//	@Router			/account/password [put]
func (uc *UserController) UpdatePassword(c echo.Context) error {
	// Bind request body to PasswordBody struct
	var p PasswordBody
	if err := c.Bind(&p); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(p); err != nil {
		return err
	}

	// Get user email from JWT
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*auth.JwtCustomClaims)
	email := claims.Subject

	var user database.User
	if err := uc.db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	// Check if the old password is correct
	hashedPassword := utils.HashPasswordMD5(user.Salt, p.OldPassword)
	if hashedPassword != user.Password {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}

	// Generate new salt and hash the new password
	salt, err := utils.GenerateRandomString(20)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	hashedNewPassword := utils.HashPasswordMD5(salt, p.NewPassword)

	// Update the user's password and salt in the database
	if err := uc.db.Model(&user).Updates(database.User{
		Salt:     salt,
		Password: hashedNewPassword,
	}).Error; err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return c.String(http.StatusOK, "password updated")
}
