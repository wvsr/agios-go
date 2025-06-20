package handlers

import (
	"net/http"

	"agios/internal/repositories"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DeleteMessageHandler(messageRepo repositories.MessageRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		messageIDStr := c.Param("messageId")
		messageID, err := uuid.Parse(messageIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": map[string]string{
					"message": "Invalid message ID format.",
					"code":    "INVALID_MESSAGE_ID",
				},
			})
		}

		err = messageRepo.DeleteMessage(c.Request().Context(), messageID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"error": map[string]string{
						"message": "Message not found.",
						"code":    "MESSAGE_NOT_FOUND",
					},
				})
			}

			c.Logger().Errorf("Error deleting message %s: %v", messageIDStr, err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]string{
					"message": "Failed to delete message.",
					"code":    "INTERNAL_ERROR",
				},
			})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
