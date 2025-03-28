package auth

import (
	"errors"
	"fmt"
	"msnserver/pkg/database"
	"msnserver/pkg/utils"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type User struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=16,alphanum"`
	FirstName string `json:"first_name" validate:"required,min=2,name"`
	LastName  string `json:"last_name" validate:"required,min=2,name"`
	Country   string `json:"country" validate:"required,country"`
	State     string `json:"state" validate:"required_if=Country US,excluded_unless=Country US,omitempty,us_state"`
	City      string `json:"city" validate:"required_if=Country US,excluded_unless=Country US,omitempty,min=2,name"`
}

// Register godoc
//
//	@Summary		Register
//	@Description	Register a new user
//	@Tags			auth
//	@Accept			json
//	@Param			user body User true "user information"
//	@Produce		plain
//	@Success		200	{string}	string	"user created"
//	@Failure		400	{string}	string	"bad request"
//	@Failure		409	{string}	string	"email already exists"
//	@Failure		500	{string}	string	"internal server error"
//	@Router			/auth/register [post]
func (ac *AuthController) Register(c echo.Context) error {
	var u User
	if err := c.Bind(&u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(u); err != nil {
		return err
	}

	email := strings.ToLower(u.Email)
	salt, err := utils.GenerateRandomString(20)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	hashedPassword := utils.HashPasswordMD5(salt, u.Password)

	firstName := utils.FormatString(u.FirstName)
	lastName := utils.FormatString(u.LastName)
	displayName := url.PathEscape(fmt.Sprintf("%s %s", firstName, lastName))
	city := utils.FormatString(u.City)

	user := database.User{
		Email:       email,
		Salt:        salt,
		Password:    hashedPassword,
		FirstName:   firstName,
		LastName:    lastName,
		Country:     u.Country,
		State:       u.State,
		City:        city,
		DisplayName: displayName,
	}

	err = ac.db.Create(&user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return c.String(http.StatusConflict, "email already exists")
	} else if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return c.String(http.StatusOK, "user created")
}
