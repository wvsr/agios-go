package handlers

import (
	"net/http"

	"agios/internal/repositories"
	"agios/internal/utils/helpers"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DeleteThreadHandler(threadRepo repositories.ThreadRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		threadIDStr := c.Param("threadId")
		threadID, err := uuid.Parse(threadIDStr)
		if err != nil {
			return helpers.JSONError(c, http.StatusBadRequest, "Invalid thread ID format.", "INVALID_THREAD_ID")
		}

		err = threadRepo.DeleteThread(c.Request().Context(), threadID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return helpers.JSONError(c, http.StatusNotFound, "Thread not found.", "THREAD_NOT_FOUND")

			}

			c.Logger().Errorf("Error deleting thread %s: %v", threadIDStr, err)

			return helpers.JSONError(c, http.StatusInternalServerError, "Failed to delete thread.", "INTERNAL_ERROR")

		}

		return c.NoContent(http.StatusNoContent)
	}
}
