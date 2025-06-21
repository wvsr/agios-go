package handlers

import (
	"net/http"

	"agios/internal/repositories"
	"agios/internal/utils/helpers"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ThreadHandler struct {
	ThreadRepo repositories.ThreadRepository
}

// @Summary Get a thread by ID
// @Description Get a thread and its messages by thread ID
// @Tags Threads
// @Accept json
// @Produce json
// @Param threadId path string true "Thread ID"
// @Success 200 {object} object{id=string,slug=string,created_at=string,updated_at=string,version=int,messages=[]object{id=string,query_text=string,response_text=string,event_type=string,stream_status=string,meta_data=string,message_index=int,created_at=string,version=int}} "Thread details with messages"
// @Failure 400 {object} helpers.ErrorResponse "Invalid thread ID format"
// @Failure 404 {object} helpers.ErrorResponse "Thread not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Router /api/v1/threads/{threadId} [get]
func (h *ThreadHandler) GetThread(c echo.Context) error {
	threadIDStr := c.Param("threadId")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		return helpers.JSONError(c, http.StatusBadRequest, "Invalid thread ID format.", "INVALID_THREAD_ID")
	}

	thread, err := h.ThreadRepo.GetThreadWithMessages(c.Request().Context(), threadID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return helpers.JSONError(c, http.StatusNotFound, "Thread not found.", "THREAD_NOT_FOUND")
		}
		return helpers.JSONError(c, http.StatusInternalServerError, "Failed to retrieve thread.", "INTERNAL_ERROR")
	}

	resp := echo.Map{
		"id":         thread.ID,
		"slug":       thread.Slug,
		"created_at": thread.CreatedAt,
		"updated_at": thread.UpdatedAt,
		"version":    thread.Version,
		"messages":   []echo.Map{},
	}

	for _, m := range thread.Messages {
		resp["messages"] = append(resp["messages"].([]echo.Map), echo.Map{
			"id":            m.ID,
			"query_text":    m.QueryText,
			"response_text": m.ResponseText,
			"event_type":    m.EventType,
			"stream_status": m.StreamStatus,
			"meta_data":     m.MetaData,
			"message_index": m.MessageIndex,
			"created_at":    m.CreatedAt,
			"version":       m.Version,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func GetThreadHandler(threadRepo repositories.ThreadRepository) echo.HandlerFunc {
	handler := &ThreadHandler{ThreadRepo: threadRepo}
	return handler.GetThread
}
