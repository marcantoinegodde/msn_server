package auth

import (
	"msnserver/config"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type AuthController struct {
	c        *config.WebServer
	db       *gorm.DB
	webauthn *webauthn.WebAuthn
}

func NewAuthController(c *config.WebServer, db *gorm.DB, webauthn *webauthn.WebAuthn) *AuthController {
	return &AuthController{
		c:        c,
		db:       db,
		webauthn: webauthn,
	}
}
