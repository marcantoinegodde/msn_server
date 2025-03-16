package web

import (
	"msnserver/config"
	"msnserver/internal/web/auth"
	"msnserver/internal/web/user"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     ws.c.WebServer.AllowedOrigins,
		AllowCredentials: true,
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
	apiGroup.GET("/healthz", Healthz)

	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/register", ac.Register)
	authGroup.POST("/login", ac.Login)

	restrictedGroup := apiGroup.Group("")
	restrictedGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(ws.c.WebServer.JWTSecret),
		TokenLookup: "header:Authorization:Bearer ,cookie:token",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
	}))

	userGroup := restrictedGroup.Group("/user")
	userGroup.GET("/me", uc.Me)

	e.Logger.Fatal(e.Start(":8080"))
}
