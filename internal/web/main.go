package web

import (
	"msnserver/config"
	"msnserver/internal/web/auth"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type WebServer struct {
	c  *config.MSNServerConfiguration
	db *gorm.DB
}

func NewWebServer(c *config.MSNServerConfiguration, db *gorm.DB) *WebServer {
	return &WebServer{
		c:  c,
		db: db,
	}
}

func (ws *WebServer) Start() {
	e := echo.New()

	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("name", validateName)
	v.RegisterValidation("country", validateCountry)
	v.RegisterValidation("us_state", validateUSState)
	e.Validator = &CustomValidator{validator: v}

	ac := auth.NewAuthController(&ws.c.WebServer, ws.db)

	e.GET("/healthz", Healthz)
	e.POST("/register", ac.Register)
	e.POST("/login", ac.Login)

	e.Logger.Fatal(e.Start(":8080"))
}
