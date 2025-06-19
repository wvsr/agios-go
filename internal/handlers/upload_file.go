package handlers

import (
	"fmt"
	"net/http"

	"agios/internal/services"

	"github.com/labstack/echo/v4"
)

const (
	maxFileSize  = 10 * 1024 * 1024 // 10MB
	maxFileCount = 5
	uploadDir    = "uploads"
)

type UploadError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type UploadResponse struct {
	FileIDs []string      `json:"fileIds"`
	Errors  []UploadError `json:"errors"`
}

func UploadFile(fileService services.FileService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the multipart form
		form, err := c.MultipartForm()
		if err != nil {
			return c.JSON(http.StatusBadRequest, UploadResponse{Errors: []UploadError{{Field: "form", Message: "Invalid form data"}}})
		}

		if form == nil {
			return c.JSON(http.StatusBadRequest, UploadResponse{Errors: []UploadError{{Field: "form", Message: "Malformed form data"}}})
		}

		files := form.File["files"]
		if len(files) == 0 {
			return c.JSON(http.StatusBadRequest, UploadResponse{Errors: []UploadError{{Field: "files", Message: "No files uploaded"}}})
		}

		if len(files) > maxFileCount {
			return c.JSON(http.StatusBadRequest, UploadResponse{Errors: []UploadError{{Field: "files", Message: fmt.Sprintf("Maximum %d files allowed", maxFileCount)}}})
		}

		var fileIDs []string
		var uploadErrors []UploadError

		for _, fileHeader := range files {
			// Check file size
			if fileHeader.Size > maxFileSize {
				uploadErrors = append(uploadErrors, UploadError{Field: fileHeader.Filename, Message: fmt.Sprintf("File size exceeds the limit of %dMB", maxFileSize/1024/1024)})
				continue
			}

			src, err := fileHeader.Open()
			if err != nil {
				uploadErrors = append(uploadErrors, UploadError{Field: fileHeader.Filename, Message: "Failed to open file"})
				continue
			}

			fileID, err := fileService.UploadFile(src, fileHeader.Filename)
			if err != nil {
				uploadErrors = append(uploadErrors, UploadError{Field: fileHeader.Filename, Message: err.Error()})
				continue
			}

			fileIDs = append(fileIDs, fileID)
		}

		response := UploadResponse{
			FileIDs: fileIDs,
			Errors:  uploadErrors,
		}

		if len(uploadErrors) > 0 && len(fileIDs) > 0 {
			return c.JSON(207, response)
		}

		if len(uploadErrors) > 0 && len(fileIDs) == 0 {
			return c.JSON(http.StatusBadRequest, response)
		}

		return c.JSON(http.StatusCreated, response)
	}
}
