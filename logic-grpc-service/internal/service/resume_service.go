package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"

	"logic-grpc-service/internal/model"
	"logic-grpc-service/internal/oss"
	"logic-grpc-service/pkg"
	"logic-grpc-service/proto/gen/resume"
)

type ResumeService struct {
	resume.UnimplementedResumeServiceServer
	db *gorm.DB
}

func NewResumeService(db *gorm.DB) *ResumeService {
	return &ResumeService{db: db}
}

func (s *ResumeService) GetUploadSignURL(ctx context.Context, req *resume.GetUploadSignURLReq) (*resume.GetUploadSignURLResp, error) {
	ossKey := oss.GenerateObjectKey(req.CandidateId, req.FileName)

	uploadURL, err := oss.GenerateUploadSignURL(ossKey, 3600, req.ContentType)
	if err != nil {
		return nil, fmt.Errorf("获取上传签名URL失败: %w", err)
	}

	return &resume.GetUploadSignURLResp{
		UploadUrl: uploadURL,
		OssKey:    ossKey,
		ExpireSec: 3600,
	}, nil
}

func (s *ResumeService) ConfirmUpload(ctx context.Context, req *resume.ConfirmUploadReq) (*resume.ConfirmUploadResp, error) {
	exists, err := oss.CheckObjectExists(req.OssKey)
	if err != nil {
		return nil, fmt.Errorf("检查文件存在性失败: %w", err)
	}
	if !exists {
		return &resume.ConfirmUploadResp{
			ResumeId: 0,
			Message:  "文件不存在，请先上传",
		}, nil
	}

	if !pkg.IsValidResumeType(req.FileType) {
		return &resume.ConfirmUploadResp{
			ResumeId: 0,
			Message:  "不支持的文件格式，仅支持 PDF、DOC、DOCX",
		}, nil
	}

	fileExt := strings.ToLower(strings.TrimPrefix(req.FileType, "."))
	if !pkg.IsValidResumeType(fileExt) {
		if err := oss.DeleteObject(req.OssKey); err != nil {
			log.Printf("删除OSS文件失败: oss_key=%s, err=%v", req.OssKey, err)
		}
		return &resume.ConfirmUploadResp{
			ResumeId: 0,
			Message:  fmt.Sprintf("文件类型与扩展名不匹配: %s", req.FileType),
		}, nil
	}

	header, err := oss.GetObjectBytes(req.OssKey, 8)
	if err != nil {
		return nil, fmt.Errorf("读取文件头失败: %w", err)
	}
	if !pkg.ValidateFileType(fileExt, header) {
		return &resume.ConfirmUploadResp{
			ResumeId: 0,
			Message:  "文件内容与格式不匹配，仅支持 PDF、DOC、DOCX 格式",
		}, nil
	}

	resumeModel := model.Resume{
		CandidateID: req.CandidateId,
		FileName:    req.FileName,
		FileType:    req.FileType,
		FileSize:    req.FileSize,
		OSSKey:      req.OssKey,
		UploadedAt:  time.Now(),
	}

	if err := s.db.Create(&resumeModel).Error; err != nil {
		return nil, fmt.Errorf("保存简历记录失败: %w", err)
	}

	return &resume.ConfirmUploadResp{
		ResumeId: resumeModel.ID,
		Message:  "上传成功",
	}, nil
}

func (s *ResumeService) GetDownloadSignURL(ctx context.Context, req *resume.GetDownloadSignURLReq) (*resume.GetDownloadSignURLResp, error) {
	var resumeModel model.Resume
	if err := s.db.Where("id = ?", req.ResumeId).First(&resumeModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &resume.GetDownloadSignURLResp{}, nil
		}
		return nil, fmt.Errorf("查询简历失败: %w", err)
	}

	if req.Role == model.RoleCandidate && resumeModel.CandidateID != req.RequesterId {
		return &resume.GetDownloadSignURLResp{}, nil
	}

	if req.Role == model.RoleHR {
		var app model.Application
		if err := s.db.Where("resume_id = ?", req.ResumeId).First(&app).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return &resume.GetDownloadSignURLResp{}, nil
			}
			return nil, fmt.Errorf("查询投递记录失败: %w", err)
		}

		var pos model.Position
		if err := s.db.Where("id = ? AND hr_id = ?", app.PositionID, req.RequesterId).First(&pos).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return &resume.GetDownloadSignURLResp{}, nil
			}
			return nil, fmt.Errorf("查询岗位失败: %w", err)
		}
	}

	downloadURL, err := oss.GenerateDownloadSignURL(resumeModel.OSSKey, 3600)
	if err != nil {
		return nil, fmt.Errorf("获取下载签名URL失败: %w", err)
	}

	return &resume.GetDownloadSignURLResp{
		DownloadUrl: downloadURL,
		FileName:    resumeModel.FileName,
		ExpireSec:   3600,
	}, nil
}

func (s *ResumeService) ListMyResumes(ctx context.Context, req *resume.ListMyResumesReq) (*resume.ListMyResumesResp, error) {
	var resumes []model.Resume
	if err := s.db.Where("candidate_id = ?", req.CandidateId).Order("uploaded_at DESC").Find(&resumes).Error; err != nil {
		return nil, fmt.Errorf("查询简历列表失败: %w", err)
	}

	result := make([]*resume.ResumeInfo, 0, len(resumes))
	for _, r := range resumes {
		result = append(result, &resume.ResumeInfo{
			Id:         r.ID,
			FileName:   r.FileName,
			FileType:   r.FileType,
			FileSize:   r.FileSize,
			UploadedAt: r.UploadedAt.Format(time.RFC3339),
		})
	}

	return &resume.ListMyResumesResp{
		Resumes: result,
	}, nil
}

func (s *ResumeService) DeleteResume(ctx context.Context, req *resume.DeleteResumeReq) (*resume.DeleteResumeResp, error) {
	var resumeModel model.Resume
	if err := s.db.Where("id = ? AND candidate_id = ?", req.ResumeId, req.CandidateId).First(&resumeModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &resume.DeleteResumeResp{Success: false, Message: "简历不存在或无权删除"}, nil
		}
		return nil, fmt.Errorf("查询简历失败: %w", err)
	}

	if err := oss.DeleteObject(resumeModel.OSSKey); err != nil {
		return nil, fmt.Errorf("删除OSS文件失败: %w", err)
	}

	if err := s.db.Where("resume_id = ?", req.ResumeId).Delete(&model.Application{}).Error; err != nil {
		return nil, fmt.Errorf("删除投递记录失败: %w", err)
	}

	if err := s.db.Delete(&resumeModel).Error; err != nil {
		return nil, fmt.Errorf("删除数据库记录失败: %w", err)
	}

	return &resume.DeleteResumeResp{Success: true, Message: "删除成功"}, nil
}
