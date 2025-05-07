package auth

import (
	"fmt"
	"log"
	"msnserver/pkg/database"
	"net/http"
	"slices"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// WebauthnLoginBegin godoc
//
//	@Summary		Begin webauthn login
//	@Description	Start the webauthn login process
//	@Tags			webauthn
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	protocol.PublicKeyCredentialRequestOptions
//	@Failure		500	{string}	string	"internal server error"
//	@Router			/auth/webauthn/login/begin [post]
func (ac *AuthController) LoginBegin(c echo.Context) error {
	// Initialize session
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	// Initialize webauthn login
	options, session, err := ac.webauthn.BeginDiscoverableLogin()
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

// WebauthnLoginFinish godoc
//
//	@Summary		Finish webauthn login
//	@Description	Finish the webauthn login process
//	@Tags			webauthn
//	@Accept			json
//	@Produce		json
//	@Param			body	body		protocol.CredentialAssertionResponse	true	"webauthn credential assertion data"
//	@Success		200		{string}	string									"login success"
//	@Failure		500		{string}	string									"internal server error"
//	@Router			/auth/webauthn/login/finish [post]
func (ac *AuthController) LoginFinish(c echo.Context) error {
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

	// Retrieve the user and validate the login
	var user database.User
	credential, err := ac.webauthn.FinishDiscoverableLogin(handler(ac.db, &user), *sessionData, c.Request())
	if err != nil {
		return err
	}

	if credential.Authenticator.CloneWarning {
		log.Println("Passkey authenticator clone warning")
	}

	// Update user's credential
	for i, cred := range user.WebauthnCredentials {
		if slices.Equal(cred.ID, credential.ID) {
			user.WebauthnCredentials[i] = *credential
			break
		}
	}

	if err := ac.db.Save(&user).Error; err != nil {
		return err
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

	secure := ac.c.Env != "development"
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

	// Clear the session data
	delete(sess.Values, "webauthn")
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "login success")
}

func handler(db *gorm.DB, user *database.User) webauthn.DiscoverableUserHandler {
	return func(rawID, userHandle []byte) (webauthn.User, error) {
		if err := db.Where("webauthn_id = ?", userHandle).First(user).Error; err != nil {
			return nil, err
		}
		return user, nil
	}
}
