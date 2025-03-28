package auth

import (
	"fmt"
	"msnserver/pkg/database"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login godoc
//
//	@Summary		Login
//	@Description	Login to the application
//	@Tags			auth
//	@Accept			json
//	@Produce		plain
//	@Param			credentials	body		LoginCredentials	true	"login credentials"
//	@Success		200			{string}	string				"login success"
//	@Failure		400			{string}	string				"bad request"
//	@Failure		401			{string}	string				"unauthorized"
//	@Failure		500			{string}	string				"internal server error"
//	@Router			/auth/login [post]
func (ac *AuthController) Login(c echo.Context) error {
	var lc LoginCredentials
	if err := c.Bind(&lc); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	var user database.User
	if err := ac.db.Where("email = ?", strings.ToLower(lc.Email)).First(&user).Error; err != nil {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}

	hashedPassword := hashPassword(user.Salt, lc.Password)
	if hashedPassword != user.Password {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}

	claims := &JwtCustomClaims{
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

	var secure bool
	if ac.c.Env == "development" {
		secure = false
	} else {
		secure = true
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    t,
		Expires:  claims.ExpiresAt.Time,
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "login success")
}
