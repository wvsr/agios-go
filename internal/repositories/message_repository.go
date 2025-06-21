package repositories

import (
	"context"

	"agios/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository interface {
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error
	CreateMessage(ctx context.Context, message *models.Message) error
}

type messageRepo struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepo{db: db}
}

func (r *messageRepo) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Message{}, messageID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *messageRepo) CreateMessage(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}
