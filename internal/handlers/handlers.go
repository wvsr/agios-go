package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler functions will be defined here
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
