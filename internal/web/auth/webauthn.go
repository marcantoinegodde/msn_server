package auth

import (
	"fmt"
	"msnserver/pkg/database"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (ac *AuthController) RegisterBegin(c echo.Context) error {
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*JwtCustomClaims)
	email := claims.Subject

	// Initialize session
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	// Retrieve user from database
	var user database.User
	if err := ac.db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	// Initialize webauthn registration
	options, session, err := ac.webauthn.BeginRegistration(&user)
	if err != nil {
		return err
	}

	// Store the webauthn session in the session store
	sess.Values["webauthn"] = session
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, options.Response)
}

func (ac *AuthController) RegisterFinish(c echo.Context) error {
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*JwtCustomClaims)
	email := claims.Subject

	// Fetch the session
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	// Retrieve the webauthn session data from the session store
	val := sess.Values["webauthn"]
	sessionData, ok := val.(*webauthn.SessionData)
	if !ok {
		return fmt.Errorf("session data not found")
	}

	// Retrieve user from database
	var user database.User
	if err := ac.db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	// Finish webauthn registration
	credential, err := ac.webauthn.FinishRegistration(&user, *sessionData, c.Request())
	if err != nil {
		return err
	}

	// Save the credential to the user
	user.WebauthnCredentials = append(user.WebauthnCredentials, *credential)
	if err := ac.db.Save(&user).Error; err != nil {
		return err
	}

	// Clear the session data
	delete(sess.Values, "webauthn")
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return c.String(http.StatusOK, "registration success")
}
