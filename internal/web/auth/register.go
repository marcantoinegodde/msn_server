package auth

import (
	"fmt"
	"msnserver/pkg/database"
	"msnserver/pkg/utils"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func (ac *AuthController) Register(c echo.Context) error {
	var u User
	if err := c.Bind(&u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(u); err != nil {
		return err
	}

	var user database.User
	if err := ac.db.Where("email = ?", u.Email).First(&user).Error; err == nil {
		return c.String(http.StatusConflict, "email already exists")
	}

	salt := utils.GenerateRandomString(20)
	hashedPassword := hashPassword(salt, u.Password)

	firstName := formatName(u.FirstName)
	lastName := formatName(u.LastName)
	displayName := url.PathEscape(fmt.Sprintf("%s %s", firstName, lastName))

	user = database.User{
		Email:       u.Email,
		Salt:        salt,
		Password:    hashedPassword,
		FirstName:   firstName,
		LastName:    lastName,
		Country:     u.Country,
		State:       u.State,
		City:        u.City,
		DisplayName: displayName,
	}

	if err := ac.db.Create(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return c.String(http.StatusOK, "user created")
}
