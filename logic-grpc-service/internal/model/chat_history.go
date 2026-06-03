package model

import (
	"time"

	"gorm.io/gorm"
)

type ChatHistory struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	HRID      int64          `gorm:"column:hr_id"`
	SessionID string         `gorm:"column:session_id"`
	Role      string         `gorm:"column:role"`
	Content   string         `gorm:"column:content"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ChatHistory) TableName() string {
	return "chat_history"
}

const (
	ChatRoleUser      = "user"
	ChatRoleAssistant = "assistant"
)