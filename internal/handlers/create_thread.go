package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateThread(c echo.Context) error {
	return c.String(http.StatusOK, "CreateThread Placeholder")
}
