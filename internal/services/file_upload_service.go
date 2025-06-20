package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"agios/internal/models"
	"agios/internal/repositories"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	uploadDir   = "uploads"
)

// UploadResult defines the response schema for an uploaded file.
type UploadResult struct {
	ID               uuid.UUID `json:"id"`
	FileName         string    `json:"file_name"`
	OriginalFileName string    `json:"original_file_name"`
	FileSizeBytes    int64     `json:"file_size_bytes"`
	MimeType         string    `json:"mime_type"`
	UploadedAt       time.Time `json:"uploaded_at"`
	Version          string    `json:"version"`
}

// FileService defines file upload operations.
type FileService interface {
	UploadSingle(ctx context.Context, fh *multipart.FileHeader) (UploadResult, error)
}

// Error definitions
var (
	ErrFileTooLarge    = fmt.Errorf("FILE_TOO_LARGE")
	ErrUnsupportedType = fmt.Errorf("UNSUPPORTED_FILE_TYPE")
)

// ErrorCode maps errors to API codes.
func ErrorCode(err error) string {
	switch err {
	case ErrFileTooLarge:
		return "FILE_TOO_LARGE"
	case ErrUnsupportedType:
		return "UNSUPPORTED_FILE_TYPE"
	default:
		return "UPLOAD_ERROR"
	}
}

// NewFileService constructs a FileService.
func NewFileService(repo repositories.FileRepository) FileService {
	return &fileServiceImpl{repo: repo}
}

type fileServiceImpl struct {
	repo repositories.FileRepository
}

// UploadSingle processes, validates, stores, and persists file metadata.
func (s *fileServiceImpl) UploadSingle(ctx context.Context, fh *multipart.FileHeader) (UploadResult, error) {
	// Validate size
	if fh.Size > maxFileSize {
		return UploadResult{}, ErrFileTooLarge
	}

	// Open source
	src, err := fh.Open()
	if err != nil {
		return UploadResult{}, err
	}
	defer src.Close()

	// Peek to detect MIME
	buf := make([]byte, 512)
	n, _ := src.Read(buf)
	mimeType := http.DetectContentType(buf[:n])

	// Validate MIME
	if !isSupported(mimeType) {
		return UploadResult{}, ErrUnsupportedType
	}

	// Reset reader
	src.Seek(0, io.SeekStart)

	// Sanitize filename
	policy := bluemonday.UGCPolicy()
	safeName := policy.Sanitize(fh.Filename)

	// Build storage path
	id := uuid.New()
	ext := filepath.Ext(safeName)
	fileName := fmt.Sprintf("%s%s", id.String(), ext)
	path := filepath.Join(uploadDir, fileName)

	// Persist
	s.repo.CreateDirectory(uploadDir)
	if err := s.repo.SaveFile(path, src); err != nil {
		return UploadResult{}, err
	}

	u := models.UploadFile{
		ID:               id,
		FileName:         fileName,
		OriginalFileName: fh.Filename,
		FileSizeBytes:    fh.Size,
		MimeType:         mimeType,
		UploadedAt:       time.Now(),
		Version:          "1.0",
	}
	if err := s.repo.SaveMetadata(ctx, &u); err != nil {
		return UploadResult{}, err
	}

	return UploadResult{
		ID:               u.ID,
		FileName:         u.FileName,
		OriginalFileName: u.OriginalFileName,
		FileSizeBytes:    u.FileSizeBytes,
		MimeType:         u.MimeType,
		UploadedAt:       u.UploadedAt,
		Version:          u.Version,
	}, nil
}

func isSupported(mime string) bool {
	supported := []string{
		"application/pdf",
		"application/x-javascript", "text/javascript",
		"application/x-python", "text/x-python",
		"text/plain", "text/html", "text/css", "text/md", "text/csv", "text/xml", "text/rtf",
		"image/png", "image/jpeg", "image/webp", "image/heic", "image/heif",
	}
	for _, m := range supported {
		if mime == m || strings.HasPrefix(mime, strings.Split(m, "/")[0]+"/") {
			return true
		}
	}
	return false
}
