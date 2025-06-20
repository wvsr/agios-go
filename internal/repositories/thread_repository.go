package repositories

import (
	"context"

	"agios/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreadRepository interface {
	DeleteThread(ctx context.Context, threadID uuid.UUID) error
}

type threadRepo struct {
	db *gorm.DB
}

func NewThreadRepository(db *gorm.DB) ThreadRepository {
	return &threadRepo{db: db}
}

func (r *threadRepo) DeleteThread(ctx context.Context, threadID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Thread{}, threadID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
