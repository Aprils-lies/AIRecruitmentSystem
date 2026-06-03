package model

import (
	"time"

	"gorm.io/gorm"
)

type Application struct {
	ID          int64          `gorm:"column:id;primaryKey"`
	PositionID  int64          `gorm:"column:position_id"`
	CandidateID int64          `gorm:"column:candidate_id"`
	ResumeID    *int64         `gorm:"column:resume_id"`
	Status      string         `gorm:"column:status"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Application) TableName() string {
	return "applications"
}

const (
	ApplicationStatusPending   = "pending"
	ApplicationStatusReviewed  = "reviewed"
	ApplicationStatusRejected  = "rejected"
	ApplicationStatusAccepted  = "accepted"
)