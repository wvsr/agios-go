package repositories

import (
	"context"

	"agios/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreadRepository interface {
	DeleteThread(ctx context.Context, threadID uuid.UUID) error
	GetThreadWithMessages(ctx context.Context, threadID uuid.UUID) (*models.Thread, error)
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

func (r *threadRepo) GetThreadWithMessages(ctx context.Context, threadID uuid.UUID) (*models.Thread, error) {
	var thread models.Thread
	result := r.db.WithContext(ctx).Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("messages.created_at ASC")
	}).Where("id = ?", threadID).First(&thread)
	if result.Error != nil {
		return nil, result.Error
	}
	return &thread, nil
}
