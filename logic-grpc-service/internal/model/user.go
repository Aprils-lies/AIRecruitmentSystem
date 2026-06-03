package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           int64          `gorm:"column:id;primaryKey"`
	Username     string         `gorm:"column:username"`
	PasswordHash string         `gorm:"column:password_hash"`
	Role         string         `gorm:"column:role"`
	RealName     *string        `gorm:"column:real_name"`
	Phone        *string        `gorm:"column:phone"`
	Education    *string        `gorm:"column:education"`
	School       *string        `gorm:"column:school"`
	Experience   *string        `gorm:"column:experience"`
	Skills       *string        `gorm:"column:skills"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (User) TableName() string {
	return "users"
}

const (
	RoleHR        = "hr"
	RoleCandidate = "candidate"
)