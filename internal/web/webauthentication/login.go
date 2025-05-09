package webauthentication

import (
	"fmt"
	"log"
	"msnserver/internal/web/auth"
	"msnserver/pkg/database"
	"net/http"
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
//	@Router			/webauthn/login/begin [post]
func (ac *WebauthnController) LoginBegin(c echo.Context) error {
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
//	@Router			/webauthn/login/finish [post]
func (ac *WebauthnController) LoginFinish(c echo.Context) error {
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

	credential.Authenticator.UpdateCounter(1)

	// Update user's credential
	cred := database.Credential{
		KeyID:              credential.ID,
		PublicKey:          credential.PublicKey,
		AttestationType:    credential.AttestationType,
		Transport:          credential.Transport,
		UserPresent:        credential.Flags.UserPresent,
		UserVerified:       credential.Flags.UserVerified,
		BackupEligible:     credential.Flags.BackupEligible,
		BackupState:        credential.Flags.BackupState,
		AAGUID:             credential.Authenticator.AAGUID,
		SignCount:          credential.Authenticator.SignCount,
		CloneWarning:       credential.Authenticator.CloneWarning,
		Attachment:         credential.Authenticator.Attachment,
		ClientDataJSON:     credential.Attestation.ClientDataJSON,
		ClientDataHash:     credential.Attestation.ClientDataHash,
		AuthenticatorData:  credential.Attestation.AuthenticatorData,
		PublicKeyAlgorithm: credential.Attestation.PublicKeyAlgorithm,
		Object:             credential.Attestation.Object,
	}
	if err := ac.db.Model(&database.Credential{}).Where("key_id = ?", cred.KeyID).Updates(cred).Error; err != nil {
		return err
	}

	claims := &auth.JwtCustomClaims{
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
		if err := db.Preload("WebauthnCredentials").Where("webauthn_id = ?", userHandle).First(user).Error; err != nil {
			return nil, err
		}
		return user, nil
	}
}
