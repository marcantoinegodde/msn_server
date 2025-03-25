package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (ac *AuthController) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "logout success")
}
