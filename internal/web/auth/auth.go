package auth

import (
	"msnserver/config"

	"gorm.io/gorm"
)

type AuthController struct {
	c  *config.WebServer
	db *gorm.DB
}

func NewAuthController(c *config.WebServer, db *gorm.DB) *AuthController {
	return &AuthController{
		c:  c,
		db: db,
	}
}
