package helpers

import (
	"github.com/labstack/echo/v4"
)

// ErrorResponse represents the error response structure.
// @model
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

func JSONError(c echo.Context, status int, message, code string) error {
	return c.JSON(status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"code":    code,
		},
	})
}
