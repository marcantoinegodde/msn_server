package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Logout godoc
//
//	@Summary		Logout
//	@Description	Logout from the application
//	@Tags			auth
//	@Accept			json
//	@Produce		plain
//	@Success		200	{string}	string	"logout success"
//	@Failure		500	{string}	string	"internal server error"
//	@Router			/auth/logout [post]
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
