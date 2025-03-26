package web

import (
	"embed"
	"msnserver/config"
	"msnserver/internal/web/auth"
	"msnserver/internal/web/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"

	_ "msnserver/internal/web/docs"
)

//go:embed all:dist
var ui embed.FS

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

//	@title			MSN Server API
//	@version		1.0
//	@description	This is the API for the MSN server web application.

//	@license.name	CC0 1.0 Universal
//	@license.url	https://creativecommons.org/publicdomain/zero/1.0/

// @BasePath	/api
func (ws *WebServer) Start() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     ws.c.WebServer.AllowedOrigins,
		AllowCredentials: true,
	}))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "dist",
		Filesystem: http.FS(ui),
	}))

	// Register custom validator
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("name", validateName)
	v.RegisterValidation("country", validateCountry)
	v.RegisterValidation("us_state", validateUSState)
	e.Validator = &CustomValidator{validator: v}

	// Register routes
	ac := auth.NewAuthController(&ws.c.WebServer, ws.db)
	uc := user.NewUserController(ws.db)

	apiGroup := e.Group("/api")
	apiGroup.GET("/swagger/*", echoSwagger.WrapHandler)
	apiGroup.GET("/healthz", Healthz)

	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/register", ac.Register)
	authGroup.POST("/login", ac.Login)
	authGroup.POST("/logout", ac.Logout)

	restrictedGroup := apiGroup.Group("")
	restrictedGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(ws.c.WebServer.JWTSecret),
		TokenLookup: "header:Authorization:Bearer ,cookie:token",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
	}))

	userGroup := restrictedGroup.Group("/user")
	userGroup.GET("/account", uc.GetAccount)
	userGroup.PATCH("/account", uc.UpdateAccount)

	e.Logger.Fatal(e.Start(":8080"))
}
