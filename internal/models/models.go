package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UploadFile struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FileName         string     `gorm:"type:text;not null"`
	OriginalFileName string     `gorm:"type:text;not null"`
	FileSizeBytes    int64      `gorm:"not null;check:file_size_bytes <= 10485760"`
	MimeType         string     `gorm:"type:text;not null"`
	UploadedAt       time.Time  `gorm:"autoCreateTime"`
	Version          string     `gorm:"type:text;not null;default:'1.0'"`
	Messages         []*Message `gorm:"many2many:message_files;constraint:OnDelete:CASCADE;"`
}

type Thread struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Slug      string     `gorm:"type:text;uniqueIndex;not null"`
	UserID    *uuid.UUID `gorm:"type:uuid;index"` // future auth support
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	Version   string     `gorm:"type:text;not null;default:'1.0'"`
	Messages  []Message  `gorm:"constraint:OnDelete:CASCADE;"`
}

type Message struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ThreadID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Thread       Thread         `gorm:"constraint:OnDelete:CASCADE;"`
	QueryText    *string        `gorm:"type:text"`
	ResponseText *string        `gorm:"type:text"`
	EventType    *string        `gorm:"type:text;check:event_type IN ('START','END','PLAN','WEB_RESULTS','MARKDOWN_ANSWER','RELATED_QUERIES','WIDGET')"`
	Model        string         `gorm:"type:text;not null"`
	InputToken   int            `gorm:"not null"`
	OutputToken  int            `gorm:"not null"`
	ResponseTime float64        `gorm:"not null"` // in seconds
	StreamStatus *string        `gorm:"type:text;check:stream_status IN ('IN_PROGRESS','DONE','FAILED')"`
	MessageIndex int            `gorm:"not null;index:,unique,composite:idx_thread_msgidx"`
	MetaData     datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	Version      string         `gorm:"type:text;not null;default:'1.0'"`
	Files        []*UploadFile  `gorm:"many2many:message_files;constraint:OnDelete:CASCADE;"`
}

type MessageFile struct {
	MessageID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	UploadFileID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
