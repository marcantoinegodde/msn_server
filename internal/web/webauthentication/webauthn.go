package webauthentication

import (
	"msnserver/config"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type WebauthnController struct {
	c        *config.WebServer
	db       *gorm.DB
	webauthn *webauthn.WebAuthn
}

func NewWebauthnController(c *config.WebServer, db *gorm.DB, webauthn *webauthn.WebAuthn) *WebauthnController {
	return &WebauthnController{
		c:        c,
		db:       db,
		webauthn: webauthn,
	}
}
