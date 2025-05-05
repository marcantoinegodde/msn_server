package web

import (
	"embed"
	"encoding/gob"
	"msnserver/config"
	"msnserver/internal/web/auth"
	"msnserver/internal/web/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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

func init() {
	gob.RegisterName("webauthn", &webauthn.SessionData{})
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

	// Setup basic middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configure CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     ws.c.WebServer.AllowedOrigins,
		AllowCredentials: true,
	}))

	// Configure session store
	secure := ws.c.WebServer.Env != "development"
	store := sessions.NewCookieStore([]byte(ws.c.WebServer.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	}
	e.Use(session.Middleware(store))

	// Serve static files
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

	wconfig := &webauthn.Config{
		RPDisplayName: "MSN Server",
		RPID:          ws.c.WebServer.RPID,
		RPOrigins:     ws.c.WebServer.RPOrigins,
	}

	webauthn, err := webauthn.New(wconfig)
	if err != nil {
		e.Logger.Fatal("Failed to create webauthn instance:", err)
	}

	jwtMiddlewareConfig := echojwt.Config{
		SigningKey:  []byte(ws.c.WebServer.JWTSecret),
		TokenLookup: "header:Authorization:Bearer ,cookie:token",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
	}

	// Register routes
	ac := auth.NewAuthController(&ws.c.WebServer, ws.db, webauthn)
	uc := user.NewUserController(ws.db)

	apiGroup := e.Group("/api")
	apiGroup.GET("/swagger/*", echoSwagger.WrapHandler)
	apiGroup.GET("/healthz", Healthz)

	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/register", ac.Register)
	authGroup.POST("/login", ac.Login)
	authGroup.POST("/logout", ac.Logout)

	webauthnGroup := authGroup.Group("/webauthn")
	webauthnGroup.POST("/login/begin", ac.LoginBegin)
	webauthnGroup.POST("/login/finish", ac.LoginFinish)
	webauthnRestrictedGroup := webauthnGroup.Group("")
	webauthnRestrictedGroup.Use(echojwt.WithConfig(jwtMiddlewareConfig))
	webauthnRestrictedGroup.POST("/register/begin", ac.RegisterBegin)
	webauthnRestrictedGroup.POST("/register/finish", ac.RegisterFinish)

	restrictedGroup := apiGroup.Group("")
	restrictedGroup.Use(echojwt.WithConfig(jwtMiddlewareConfig))

	userGroup := restrictedGroup.Group("/user")
	userGroup.GET("/account", uc.GetAccount)
	userGroup.PATCH("/account", uc.UpdateAccount)
	userGroup.PUT("/account/password", uc.UpdatePassword)

	e.Logger.Fatal(e.Start(":8080"))
}
