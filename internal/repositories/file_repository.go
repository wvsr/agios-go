package repositories

import (
	"context"
	"io"
	"os"

	"agios/internal/models"

	"gorm.io/gorm"
)

// FileRepository defines persistence methods.
type FileRepository interface {
	SaveFile(path string, src io.Reader) error
	CreateDirectory(dir string) error
	SaveMetadata(ctx context.Context, uf *models.UploadFile) error
}

// fileRepo implements FileRepository.
type fileRepo struct {
	db *gorm.DB
}

// NewFileRepository creates a repository with DB connection.
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepo{db: db}
}

func (r *fileRepo) SaveFile(path string, src io.Reader) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func (r *fileRepo) CreateDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func (r *fileRepo) SaveMetadata(ctx context.Context, uf *models.UploadFile) error {
	return r.db.WithContext(ctx).Create(uf).Error
}
