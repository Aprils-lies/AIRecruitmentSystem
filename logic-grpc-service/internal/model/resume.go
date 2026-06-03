package model

import (
	"time"

	"gorm.io/gorm"
)

type Resume struct {
	ID          int64          `gorm:"column:id;primaryKey"`
	CandidateID int64          `gorm:"column:candidate_id"`
	FileName    string         `gorm:"column:file_name"`
	FileType    string         `gorm:"column:file_type"`
	FileSize    int64          `gorm:"column:file_size"`
	OSSKey      string         `gorm:"column:oss_key"`
	UploadedAt  time.Time      `gorm:"column:uploaded_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Resume) TableName() string {
	return "resumes"
}

const (
	ResumeFileTypePDF  = "pdf"
	ResumeFileTypeDOC  = "doc"
	ResumeFileTypeDOCX = "docx"
)