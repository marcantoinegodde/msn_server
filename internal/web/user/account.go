package user

import (
	"msnserver/internal/web/auth"
	"msnserver/pkg/database"
	"msnserver/pkg/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserBody struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,name"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,name"`
	Country   string `json:"country" validate:"omitempty,country"`
	State     string `json:"state" validate:"required_if=Country US,excluded_unless=Country US,omitempty,us_state"`
	City      string `json:"city" validate:"required_if=Country US,excluded_unless=Country US"`
}

type UserResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Country   string `json:"country"`
	State     string `json:"state"`
	City      string `json:"city"`
}

// GetAccount godoc
//
//	@Summary		Get account
//	@Description	Get user account information
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	UserResponse
//	@Failure		500	{string}	string	"internal server error"
//	@Router			/user/account [get]
func (ac *UserController) GetAccount(c echo.Context) error {
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

// UpdateAccount godoc
//
//	@Summary		Update account
//	@Description	Update user account information
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		UserBody	true	"User information"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{string}	string	"bad request"
//	@Failure		500		{string}	string	"internal server error"
//	@Router			/user/account [patch]
func (ac *UserController) UpdateAccount(c echo.Context) error {
	// Bind request body to UserBody struct
	var u UserBody
	if err := c.Bind(&u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(u); err != nil {
		return err
	}

	// Get user email from JWT
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*auth.JwtCustomClaims)
	email := claims.Subject

	// Retrieve user from database
	var user database.User
	if err := ac.db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	firstName := utils.FormatString(u.FirstName)
	lastName := utils.FormatString(u.LastName)
	city := utils.FormatString(u.City)

	// Update user data
	updates := database.User{
		FirstName: firstName,
		LastName:  lastName,
		Country:   u.Country,
		State:     u.State,
		City:      city,
	}

	if err := ac.db.Model(&user).Updates(updates).Error; err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	// If country is not US, update state and city to null
	if updates.Country != "" && updates.Country != "US" {
		if err := ac.db.Model(&user).Select("state", "city").Updates(updates).Error; err != nil {
			return c.String(http.StatusInternalServerError, "internal server error")
		}
	}

	// Refresh user data to get updated values
	if err := ac.db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	ur := UserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Country:   user.Country,
		State:     user.State,
		City:      user.City,
	}

	return c.JSON(http.StatusOK, ur)
}
