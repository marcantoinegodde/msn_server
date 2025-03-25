package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Healthz godoc
//
//	@Summary		Healthz route
//	@Description	Get the health status of the application
//	@Tags			misc
//	@Produce		plain
//	@Success		200	{string}	string	"OK"
//	@Failure		500	{string}	string	"internal server error"
//	@Router			/healthz [get]
func Healthz(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
