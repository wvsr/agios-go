package handlers

import (
	"net/http"

	"agios/internal/repositories"
	"agios/internal/utils/helpers"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DeleteMessageHandler(messageRepo repositories.MessageRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		messageIDStr := c.Param("messageId")
		messageID, err := uuid.Parse(messageIDStr)
		if err != nil {
			return helpers.JSONError(c, http.StatusBadRequest, "Invalid message ID format", "INVALID_MESSAGE_ID")
		}

		err = messageRepo.DeleteMessage(c.Request().Context(), messageID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return helpers.JSONError(c, http.StatusNotFound, "Message not found.", "MESSAGE_NOT_FOUND")
			}

			c.Logger().Errorf("Error deleting message %s: %v", messageIDStr, err)

			return helpers.JSONError(c, http.StatusInternalServerError, "Failed to delete message.", "INTERNAL_ERROR")
		}

		return c.NoContent(http.StatusNoContent)
	}
}
