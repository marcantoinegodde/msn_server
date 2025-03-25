package user

import (
	"msnserver/internal/web/auth"
	"msnserver/pkg/database"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Country   string `json:"country"`
	State     string `json:"state"`
	City      string `json:"city"`
}

// Me godoc
//
//	@Summary		Me route
//	@Description	Get the user information
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	UserResponse
//	@Failure		500	{string}	string	"internal server error"
//	@Router			/user/me [get]
func (ac *UserController) Me(c echo.Context) error {
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*auth.JwtCustomClaims)
	email := claims.Subject

	var user database.User
	if err := ac.db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	u := UserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Country:   user.Country,
		State:     user.State,
		City:      user.City,
	}

	return c.JSON(http.StatusOK, u)
}
