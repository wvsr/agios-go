package handlers

import (
	"net/http"

	"agios/internal/repositories"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DeleteThreadHandler(threadRepo repositories.ThreadRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		threadIDStr := c.Param("threadId")
		threadID, err := uuid.Parse(threadIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": map[string]string{
					"message": "Invalid thread ID format.",
					"code":    "INVALID_THREAD_ID",
				},
			})
		}

		err = threadRepo.DeleteThread(c.Request().Context(), threadID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"error": map[string]string{
						"message": "Thread not found.",
						"code":    "THREAD_NOT_FOUND",
					},
				})
			}

			c.Logger().Errorf("Error deleting thread %s: %v", threadIDStr, err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]string{
					"message": "Failed to delete thread.",
					"code":    "INTERNAL_ERROR",
				},
			})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
