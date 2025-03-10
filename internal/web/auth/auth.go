package auth

import (
	"msnserver/config"

	"gorm.io/gorm"
)

type User struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=16,alphanum"`
	FirstName string `json:"first_name" validate:"required,min=2,name"`
	LastName  string `json:"last_name" validate:"required,min=2,name"`
	Country   string `json:"country" validate:"required,country"`
	State     string `json:"state" validate:"required_if=Country US,excluded_unless=Country US,omitempty,us_state"`
	City      string `json:"city" validate:"required_if=Country US,excluded_unless=Country US"`
}

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
