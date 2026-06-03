package model

import (
	"time"

	"gorm.io/gorm"
)

type Position struct {
	ID          int64          `gorm:"column:id;primaryKey"`
	HRID        int64          `gorm:"column:hr_id"`
	Title       string         `gorm:"column:title"`
	Description string         `gorm:"column:description"`
	Requirements string         `gorm:"column:requirements"`
	SalaryMin   *int32         `gorm:"column:salary_min"`
	SalaryMax   *int32         `gorm:"column:salary_max"`
	Location    *string        `gorm:"column:location"`
	Status      string         `gorm:"column:status"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Position) TableName() string {
	return "positions"
}

const (
	PositionStatusPublished = "published"
	PositionStatusOffline   = "offline"
)