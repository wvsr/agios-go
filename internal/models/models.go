package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Thread struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID      `gorm:"type:uuid"` // Foreign key to a users table
	Title     string         `gorm:"type:varchar(255);not null"`
	URLSlug   string         `gorm:"type:varchar(255);unique;not null"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type MessageRole string

const (
	UserRole      MessageRole = "user"
	AssistantRole MessageRole = "assistant"
)

type Message struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ThreadID       uuid.UUID      `gorm:"type:uuid;not null"`
	Thread         Thread         `gorm:"foreignKey:ThreadID;references:ID;onDelete:CASCADE"`
	Role           MessageRole    `gorm:"type:message_role;not null"`
	QueryText      string         `gorm:"type:text"`  // For 'user' roles, this contains the user'styped query.
	ResponseBlocks []byte         `gorm:"type:jsonb"` // For 'assistant' roles, this stores the final, complete array of Blocks as JSON.
	CreatedAt      time.Time      `gorm:"not null;default:now()"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type File struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MessageID        uuid.UUID      `gorm:"type:uuid;not null"`
	Message          Message        `gorm:"foreignKey:MessageID;references:ID;onDelete:CASCADE"`
	OriginalFilename string         `gorm:"type:varchar(255);not null"`
	StoragePath      string         `gorm:"type:varchar(1024);not null"`
	FileSizeByte     int64          `gorm:"not null"`
	MIMEType         string         `gorm:"type:varchar(100);not null"`
	CreatedAt        time.Time      `gorm:"not null;default:now()"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type BlockType string

const (
	PlanBlockType               BlockType = "PLAN"
	WebResultsBlockType         BlockType = "WEB_RESULTS"
	MarkdownAnswerBlockType     BlockType = "MARKDOWN_ANSWER"
	RelatedQueriesBlockType     BlockType = "RELATED_QUERIES"
	WeatherWidgetBlockType      BlockType = "WEATHER_WIDGET"
	FinanceWidgetBlockType      BlockType = "FINANCE_WIDGET"
	SynthesizerResultsBlockType BlockType = "SYNTHESIZER_RESULTS"
)

type Block struct {
	BlockId   string         `gorm:"type:varchar(255);primaryKey"` // Unique ID for this block instance (e.g., 'plan-msg2')
	BlockType BlockType      `gorm:"type:varchar(255)"`            // Enum determining the UI component to render
	Status    string         `gorm:"type:varchar(255)"`
	Payload   []byte         `gorm:"type:jsonb"` // Data structure specific to the blockType
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
