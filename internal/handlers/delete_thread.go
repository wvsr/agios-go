package handlers

import (
	"net/http"

	"agios/internal/repositories"
	"agios/internal/utils/helpers"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// @Summary Delete a thread by ID
// @Description Delete a thread by thread ID
// @Tags Threads
// @Accept json
// @Produce json
// @Param threadId path string true "Thread ID"
// @Success 204 "Thread successfully deleted"
// @Failure 400 {object} helpers.ErrorResponse "Invalid thread ID format"
// @Failure 404 {object} helpers.ErrorResponse "Thread not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Router /api/v1/threads/{threadId} [delete]
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
