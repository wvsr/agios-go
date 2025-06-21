package handlers

import (
	"agios/internal/services"
	"agios/internal/utils/helpers" // Import the helpers package
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary Upload files
// @Description Upload one or more files
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param files formData []file true "Files to upload" collectionFormat(multi)
// @Success 200 {array} services.UploadResult "Successfully uploaded files"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request or file upload failed"
// @Router /api/v1/files/upload [post]
func UploadFileHandler(fileService services.FileService) echo.HandlerFunc {
	return func(c echo.Context) error {

		err := c.Request().ParseMultipartForm(32 << 20)

		if err != nil {
			return helpers.JSONError(c, http.StatusBadRequest, "Invalid multipart form.", "INVALID_FORM")
		}

		files := c.Request().MultipartForm.File["files"]
		if len(files) == 0 {
			return helpers.JSONError(c, http.StatusBadRequest, "No files uploaded.", "NO_FILES")
		}
		if len(files) > 5 {
			return helpers.JSONError(c, http.StatusBadRequest, "Maximum 5 files allowed per upload.", "MAX_FILE_COUNT_EXCEEDED")
		}

		var results []services.UploadResult
		for _, fh := range files {
			res, err := fileService.UploadSingle(c.Request().Context(), fh)
			if err != nil {
				return helpers.JSONError(c, http.StatusBadRequest, err.Error(), services.ErrorCode(err))
			}
			results = append(results, res)
		}

		return c.JSON(http.StatusOK, results)
	}
}
