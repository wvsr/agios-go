package handlers

import (
	"net/http"

	"agios/internal/repositories"
	"agios/internal/utils/helpers"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// @Summary Delete a message by ID
// @Description Delete a message by message ID
// @Tags Messages
// @Accept json
// @Produce json
// @Param messageId path string true "Message ID"
// @Success 204 "Message successfully deleted"
// @Failure 400 {object} helpers.ErrorResponse "Invalid message ID format"
// @Failure 404 {object} helpers.ErrorResponse "Message not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Router /api/v1/messages/{messageId} [delete]
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
