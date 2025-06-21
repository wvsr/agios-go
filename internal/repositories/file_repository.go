package repositories

import (
	"context"
	"io"
	"os"

	"agios/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileRepository interface {
	SaveFile(path string, src io.Reader) error
	CreateDirectory(dir string) error
	SaveMetadata(ctx context.Context, uf *models.UploadFile) error
	GetFilesByIDs(ctx context.Context, fileIDs []string) ([]*models.UploadFile, error)
}

type fileRepo struct {
	db *gorm.DB
}

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

func (r *fileRepo) GetFilesByIDs(ctx context.Context, fileIDs []string) ([]*models.UploadFile, error) {
	var files []*models.UploadFile
	uuids := make([]uuid.UUID, len(fileIDs))
	for i, id := range fileIDs {
		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		uuids[i] = parsedUUID
	}

	result := r.db.WithContext(ctx).Find(&files, uuids)
	if result.Error != nil {
		return nil, result.Error
	}
	return files, nil
}
