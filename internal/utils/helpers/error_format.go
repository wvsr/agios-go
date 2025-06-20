package helpers

import (
	"github.com/labstack/echo/v4"
)

func JSONError(c echo.Context, status int, message, code string) error {
	return c.JSON(status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"code":    code,
		},
	})
}
