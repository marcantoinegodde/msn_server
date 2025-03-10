package auth

import (
	"fmt"
	"msnserver/pkg/database"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func (ac *AuthController) Login(c echo.Context) error {
	var u User
	if err := c.Bind(&u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	var user database.User
	if err := ac.db.Where("email = ?", u.Email).First(&user).Error; err != nil {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}

	hashedPassword := hashPassword(user.Salt, u.Password)
	if hashedPassword != user.Password {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}

	claims := &jwtCustomClaims{
		Name: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(ac.c.JWTSecret))
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
