package handlers

import (
	"agios/internal/config"
	"agios/internal/models"
	"agios/internal/repositories"
	"agios/internal/utils/constant"
	"agios/internal/utils/helpers"
	"agios/internal/utils/sse"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
)

type CreateThreadRequest struct {
	Slug      string   `json:"slug"`
	QueryText string   `json:"query_text"`
	FileIDs   []string `json:"file_ids"`
}

func CreateThreadHandler(threadRepo repositories.ThreadRepository, messageRepo repositories.MessageRepository, fileRepo repositories.FileRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		sse, err := sse.SetupSSE(c)
		if err != nil {
			return err
		}
		req := new(CreateThreadRequest)

		if err := c.Bind(req); err != nil {
			return helpers.JSONError(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		}

		if req.Slug == "" {
			return helpers.JSONError(c, http.StatusBadRequest, "slug cannot be blank", "SLUG_BLANK")
		}

		if req.QueryText == "" {
			return helpers.JSONError(c, http.StatusBadRequest, "query_text cannot be blank", "QUERY_TEXT_BLANK")
		}

		if len(req.FileIDs) > 5 {
			return helpers.JSONError(c, http.StatusBadRequest, "Maximum 5 file_ids allowed", "MAX_FILE_COUNT_EXCEEDED")
		}

		if helpers.WordCount(req.QueryText) > 1000 {
			return helpers.JSONError(c, http.StatusBadRequest, "Query text exceeds the 1000-word limit.", "QUERY_TEXT_TOO_LONG")
		}

		if !helpers.ValidateSlugFormat(req.Slug, req.QueryText) {
			return helpers.JSONError(c, http.StatusBadRequest, "Invalid slug format", "INVALID_SLUG_FORMAT")
		}

		existingThread, err := threadRepo.GetThreadBySlug(c.Request().Context(), req.Slug)
		if err != nil {
			return helpers.JSONError(c, http.StatusInternalServerError, "Database error checking for existing slug", "DATABASE_ERROR")
		}

		if existingThread != nil {
			return helpers.JSONError(c, http.StatusConflict, "Thread already exists", "SLUG_ALREADY_EXISTS")
		}

		newThread := &models.Thread{
			Slug: req.Slug,
		}

		if err := threadRepo.CreateThread(c.Request().Context(), newThread); err != nil {
			return helpers.JSONError(c, http.StatusInternalServerError, "Database error creating thread", "THREAD_CREATION_FAILED")
		}

		cfg, _ := config.LoadConfig()

		initialMessage := &models.Message{
			ThreadID:     newThread.ID,
			QueryText:    &req.QueryText,
			MessageIndex: 0,
			Model:        cfg.CurrentLLMModel,
			InputToken:   0,                                      // TODO: Calculate input tokens
			OutputToken:  0,                                      // TODO: Calculate output tokens
			ResponseTime: 0,                                      // TODO: Measure response time
			StreamStatus: helpers.StringPtr("IN_PROGRESS"),       // Initial status
			EventType:    helpers.StringPtr(constant.EventStart), // Initial event type
			MetaData:     datatypes.JSON([]byte("{}")),           // Initial empty metadata as JSON byte slice
		}

		if len(req.FileIDs) > 0 {
			files, err := fileRepo.GetFilesByIDs(c.Request().Context(), req.FileIDs)
			if err != nil {
				return helpers.JSONError(c, http.StatusInternalServerError, "Database error retrieving files", "FILE_RETRIEVAL_FAILED")
			}
			initialMessage.Files = files
			// todo: skip file if it doesnt exist rather than returning error
		}

		if err := messageRepo.CreateMessage(c.Request().Context(), initialMessage); err != nil {
			deleteErr := threadRepo.DeleteThread(c.Request().Context(), newThread.ID)
			if deleteErr != nil {
				c.Logger().Errorf("Failed to delete thread %s after message creation failed: %v", newThread.ID, deleteErr)
			}
			return helpers.JSONError(c, http.StatusInternalServerError, "Database error creating initial message", "MESSAGE_CREATION_FAILED")
		}

		// --- Send START event (already done above validation) ---
		// startEventData, _ := json.Marshal(map[string]bool{"streaming": true})
		// sse.SendEvent(constant.EventStart, string(startEventData))

		// --- Implement message processing logic here ---
		// This is where the main processing (calling services/prompts/LLMs) would happen.
		// This part is complex and likely involves background processing or goroutines
		// to keep the SSE stream open while processing occurs.
		// The SSE stream should send PLAN, WEB_RESULTS, MARKDOWN_ANSWER, WIDGET events during this phase.

		// Example of sending a PLAN event (placeholder):
		// planEventData, _ := json.Marshal(map[string]interface{}{"version": "1.0", "cot": "Planning the response...", "streaming": true})
		// sse.SendEvent(constant.EventPlan, string(planEventData))

		// TODO: Implement the actual message processing logic here.
		// This will involve:
		// - Calling LLMs/services based on req.QueryText and associated files.
		// - Updating the initialMessage in the database with response_text, event_type, stream_status, tokens, response_time, metadata.
		// - Sending subsequent SSE events (PLAN, WEB_RESULTS, MARKDOWN_ANSWER, WIDGET) via the 'sse' writer.

		// --- Send END event when processing is complete ---
		// This should be sent *after* all processing and other events are done.
		endEventData, _ := json.Marshal(map[string]bool{"streaming": false})
		sse.SendEvent(constant.EventEnd, string(endEventData))

		// Note: The function should return after the SSE stream is complete.
		// The SSE setup handles keeping the connection open.
		// Returning nil or a success status here might close the connection prematurely depending on Echo's SSE handling.
		// A common pattern is to block until the stream is done or manage the connection lifecycle explicitly.
		// For this starting point, we'll rely on the defer sse.Close() and assume the SSE library manages the connection until events are sent.
		// A more robust implementation might involve goroutines and channels to manage the async nature of SSE and background processing.

		// Returning nil indicates successful handling of the request, even though the response is streamed.
		return nil
	}
}
