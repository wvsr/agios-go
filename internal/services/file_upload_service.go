package services

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"agios/internal/repositories"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	uploadDir   = "uploads"
)

type FileService interface {
	UploadFile(file io.ReadCloser, filename string) (string, error)
}

type FileServiceImpl struct {
	fileRepository repositories.FileRepository
}

func NewFileService(fileRepository repositories.FileRepository) FileService {
	return &FileServiceImpl{
		fileRepository: fileRepository,
	}
}

func (s *FileServiceImpl) UploadFile(file io.ReadCloser, filename string) (string, error) {
	defer file.Close()

	// Check file type
	mimeType, err := detectMimeType(file)
	if err != nil {
		return "", fmt.Errorf("failed to detect file type: %w", err)
	}

	allowed := false
	if strings.HasPrefix(mimeType, "text/") ||
		strings.HasPrefix(mimeType, "image/") ||
		mimeType == "application/pdf" {
		allowed = true
	}

	if !allowed {
		return "", fmt.Errorf("invalid file type: %s", mimeType)
	}

	// Sanitize filename
	p := bluemonday.UGCPolicy()
	safeFilename := p.Sanitize(filename)

	// Generate a unique file ID
	fileID := uuid.New().String()
	dstPath := filepath.Join(uploadDir, fileID+filepath.Ext(safeFilename))

	// Save the file to disk
	err = s.fileRepository.CreateDirectory(uploadDir)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	err = s.fileRepository.SaveFile(dstPath, file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return fileID, nil
}

func detectMimeType(file io.Reader) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	return contentType, nil
}
