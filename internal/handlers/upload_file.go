package handlers

import (
	"agios/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// UploadFileHandler handles multipart file uploads with validation and metadata persistence.
func UploadFileHandler(fileService services.FileService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse multipart form
		err := c.Request().ParseMultipartForm(32 << 20)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": echo.Map{"message": "Invalid multipart form.", "code": "INVALID_FORM"},
			})
		}

		files := c.Request().MultipartForm.File["files"]
		if len(files) == 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": echo.Map{"message": "No files uploaded.", "code": "NO_FILES"}})
		}
		if len(files) > 5 {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": echo.Map{"message": "Maximum 5 files allowed per upload.", "code": "MAX_FILE_COUNT_EXCEEDED"}})
		}

		var results []services.UploadResult
		for _, fh := range files {
			res, err := fileService.UploadSingle(c.Request().Context(), fh)
			if err != nil {
				// Return specific error
				return c.JSON(http.StatusBadRequest, echo.Map{"error": echo.Map{"message": err.Error(), "code": services.ErrorCode(err)}})
			}
			results = append(results, res)
		}

		return c.JSON(http.StatusOK, results)
	}
}
