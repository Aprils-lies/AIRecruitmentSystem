package service

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"logic-grpc-service/internal/model"
	"logic-grpc-service/proto/gen/application"
)

type ApplicationService struct {
	application.UnimplementedApplicationServiceServer
	db *gorm.DB
}

func NewApplicationService(db *gorm.DB) *ApplicationService {
	return &ApplicationService{db: db}
}

func (s *ApplicationService) Apply(ctx context.Context, req *application.ApplyReq) (*application.ApplyResp, error) {
	var pos model.Position
	if err := s.db.Where("id = ? AND status = ?", req.PositionId, model.PositionStatusPublished).First(&pos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &application.ApplyResp{
				ApplicationId: 0,
				Message:       "岗位不存在或已下架",
			}, nil
		}
		return nil, fmt.Errorf("查询岗位失败: %w", err)
	}

	var candidate model.User
	if err := s.db.Where("id = ? AND role = ?", req.CandidateId, model.RoleCandidate).First(&candidate).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &application.ApplyResp{
				ApplicationId: 0,
				Message:       "候选人不存在或角色错误",
			}, nil
		}
		return nil, fmt.Errorf("查询候选人信息失败: %w", err)
	}

	if candidate.RealName == nil || *candidate.RealName == "" {
		return &application.ApplyResp{
			ApplicationId: 0,
			Message:       "请先完善个人资料：姓名不能为空",
		}, nil
	}
	if candidate.Phone == nil || *candidate.Phone == "" {
		return &application.ApplyResp{
			ApplicationId: 0,
			Message:       "请先完善个人资料：联系电话不能为空",
		}, nil
	}
	if candidate.Education == nil || *candidate.Education == "" {
		return &application.ApplyResp{
			ApplicationId: 0,
			Message:       "请先完善个人资料：学历不能为空",
		}, nil
	}
	if candidate.School == nil || *candidate.School == "" {
		return &application.ApplyResp{
			ApplicationId: 0,
			Message:       "请先完善个人资料：毕业院校不能为空",
		}, nil
	}
	if candidate.Experience == nil || *candidate.Experience == "" {
		return &application.ApplyResp{
			ApplicationId: 0,
			Message:       "请先完善个人资料：工作经历不能为空",
		}, nil
	}
	if candidate.Skills == nil || *candidate.Skills == "" {
		return &application.ApplyResp{
			ApplicationId: 0,
			Message:       "请先完善个人资料：核心技能不能为空",
		}, nil
	}

	var resumeModel model.Resume
	if err := s.db.Where("id = ? AND candidate_id = ? AND deleted_at IS NULL", req.ResumeId, req.CandidateId).First(&resumeModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &application.ApplyResp{
				ApplicationId: 0,
				Message:       "请先上传合规简历",
			}, nil
		}
		return nil, fmt.Errorf("查询简历失败: %w", err)
	}

	var existingApp model.Application
	if err := s.db.Where("position_id = ? AND candidate_id = ? AND deleted_at IS NULL", req.PositionId, req.CandidateId).First(&existingApp).Error; err == nil {
		return &application.ApplyResp{
			ApplicationId: 0,
			Message:       "您已投递过该岗位",
		}, nil
	}

	app := model.Application{
		PositionID:  req.PositionId,
		CandidateID: req.CandidateId,
		ResumeID:    &req.ResumeId,
		Status:      model.ApplicationStatusPending,
	}

	if err := s.db.Create(&app).Error; err != nil {
		return nil, fmt.Errorf("创建投递记录失败: %w", err)
	}

	return &application.ApplyResp{
		ApplicationId: app.ID,
		Message:       "投递成功",
	}, nil
}

func (s *ApplicationService) ListCandidates(ctx context.Context, req *application.ListCandidatesReq) (*application.ListCandidatesResp, error) {
	var pos model.Position
	if err := s.db.Where("id = ? AND hr_id = ?", req.PositionId, req.HrId).First(&pos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &application.ListCandidatesResp{
				Candidates: nil,
				Total:      0,
				Page:       req.Page,
			}, nil
		}
		return nil, fmt.Errorf("查询岗位失败: %w", err)
	}

	query := s.db.Model(&model.Application{}).Where("position_id = ?", req.PositionId)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("统计候选人数量失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	var applications []model.Application
	if err := query.Order("created_at DESC").Offset(int(offset)).Limit(int(req.PageSize)).Find(&applications).Error; err != nil {
		return nil, fmt.Errorf("查询投递记录失败: %w", err)
	}

	candidateIDs := make([]int64, 0, len(applications))
	for _, app := range applications {
		candidateIDs = append(candidateIDs, app.CandidateID)
	}

	var users []model.User
	if len(candidateIDs) > 0 {
		if err := s.db.Where("id IN ?", candidateIDs).Find(&users).Error; err != nil {
			return nil, fmt.Errorf("查询用户信息失败: %w", err)
		}
	}

	userMap := make(map[int64]model.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	var resumeIDs []int64
	for _, app := range applications {
		if app.ResumeID != nil {
			resumeIDs = append(resumeIDs, *app.ResumeID)
		}
	}

	var resumes []model.Resume
	if len(resumeIDs) > 0 {
		if err := s.db.Where("id IN ?", resumeIDs).Find(&resumes).Error; err != nil {
			return nil, fmt.Errorf("查询简历信息失败: %w", err)
		}
	}

	resumeMap := make(map[int64]model.Resume)
	for _, r := range resumes {
		resumeMap[r.ID] = r
	}

	result := make([]*application.CandidateBrief, 0, len(applications))
	for _, app := range applications {
		user := userMap[app.CandidateID]
		resume := resumeMap[app.CandidateID]

		result = append(result, &application.CandidateBrief{
			UserId:        user.ID,
			Username:      user.Username,
			RealName:      getStringPtr(user.RealName),
			Phone:         getStringPtr(user.Phone),
			Education:     getStringPtr(user.Education),
			School:        getStringPtr(user.School),
			Skills:        getStringPtr(user.Skills),
			ResumeId:      resume.ID,
			ResumeName:    resume.FileName,
			AppliedAt:     app.CreatedAt.Format(time.RFC3339),
			Status:        app.Status,
			ApplicationId: app.ID,
		})
	}

	return &application.ListCandidatesResp{
		Candidates: result,
		Total:      int32(total),
		Page:       req.Page,
	}, nil
}

func (s *ApplicationService) GetCandidateDetail(ctx context.Context, req *application.GetCandidateDetailReq) (*application.GetCandidateDetailResp, error) {
	var app model.Application
	if err := s.db.Where("id = ? AND candidate_id = ?", req.ApplicationId, req.CandidateId).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &application.GetCandidateDetailResp{}, nil
		}
		return nil, fmt.Errorf("查询投递记录失败: %w", err)
	}

	var pos model.Position
	if err := s.db.Where("id = ? AND hr_id = ?", app.PositionID, req.HrId).First(&pos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &application.GetCandidateDetailResp{}, nil
		}
		return nil, fmt.Errorf("查询岗位失败: %w", err)
	}

	var user model.User
	if err := s.db.Where("id = ?", req.CandidateId).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &application.GetCandidateDetailResp{}, nil
		}
		return nil, fmt.Errorf("查询用户信息失败: %w", err)
	}

	var resume model.Resume
	if app.ResumeID != nil {
		if err := s.db.Where("id = ?", *app.ResumeID).First(&resume).Error; err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("查询简历信息失败: %w", err)
		}
	}

	return &application.GetCandidateDetailResp{
		UserId:     user.ID,
		Username:   user.Username,
		RealName:   getStringPtr(user.RealName),
		Phone:      getStringPtr(user.Phone),
		Education:  getStringPtr(user.Education),
		School:     getStringPtr(user.School),
		Experience: getStringPtr(user.Experience),
		Skills:     getStringPtr(user.Skills),
		ResumeId:   resume.ID,
		ResumeName: resume.FileName,
		ResumeType: resume.FileType,
		AppliedAt:  app.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *ApplicationService) ListMyApplications(ctx context.Context, req *application.ListMyApplicationsReq) (*application.ListMyApplicationsResp, error) {
	query := s.db.Model(&model.Application{}).Where("candidate_id = ?", req.CandidateId)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("统计投递记录数量失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	var applications []model.Application
	if err := query.Order("created_at DESC").Offset(int(offset)).Limit(int(req.PageSize)).Find(&applications).Error; err != nil {
		return nil, fmt.Errorf("查询投递记录失败: %w", err)
	}

	positionIDs := make([]int64, 0, len(applications))
	for _, app := range applications {
		positionIDs = append(positionIDs, app.PositionID)
	}

	var positions []model.Position
	if len(positionIDs) > 0 {
		if err := s.db.Where("id IN ?", positionIDs).Find(&positions).Error; err != nil {
			return nil, fmt.Errorf("查询岗位信息失败: %w", err)
		}
	}

	positionMap := make(map[int64]model.Position)
	for _, pos := range positions {
		positionMap[pos.ID] = pos
	}

	result := make([]*application.MyApplication, 0, len(applications))
	for _, app := range applications {
		pos := positionMap[app.PositionID]
		result = append(result, &application.MyApplication{
			ApplicationId: app.ID,
			PositionId:    app.PositionID,
			PositionTitle: pos.Title,
			Status:        app.Status,
			AppliedAt:     app.CreatedAt.Format(time.RFC3339),
		})
	}

	return &application.ListMyApplicationsResp{
		Applications: result,
		Total:        int32(total),
		Page:         req.Page,
	}, nil
}

func (s *ApplicationService) WithdrawApplication(ctx context.Context, req *application.WithdrawApplicationReq) (*application.WithdrawApplicationResp, error) {
	var app model.Application
	if err := s.db.Where("id = ?", req.ApplicationId).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &application.WithdrawApplicationResp{
				Success: false,
				Message: "投递记录不存在",
			}, nil
		}
		return nil, fmt.Errorf("查询投递记录失败: %w", err)
	}

	if app.CandidateID != req.CandidateId {
		return &application.WithdrawApplicationResp{
			Success: false,
			Message: "只能撤回自己的投递",
		}, nil
	}

	if app.Status != model.ApplicationStatusPending {
		return &application.WithdrawApplicationResp{
			Success: false,
			Message: "只能撤回待审核的投递",
		}, nil
	}

	if err := s.db.Delete(&app).Error; err != nil {
		return nil, fmt.Errorf("撤回投递失败: %w", err)
	}

	return &application.WithdrawApplicationResp{
		Success: true,
		Message: "撤回成功",
	}, nil
}

func (s *ApplicationService) UpdateApplicationStatus(ctx context.Context, req *application.UpdateApplicationStatusReq) (*application.UpdateApplicationStatusResp, error) {
	var app model.Application
	if err := s.db.Where("id = ?", req.ApplicationId).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &application.UpdateApplicationStatusResp{
				Success: false,
				Message: "投递记录不存在",
			}, nil
		}
		return nil, fmt.Errorf("查询投递记录失败: %w", err)
	}

	var pos model.Position
	if err := s.db.Where("id = ? AND hr_id = ?", app.PositionID, req.HrId).First(&pos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &application.UpdateApplicationStatusResp{
				Success: false,
				Message: "无权操作该投递记录",
			}, nil
		}
		return nil, fmt.Errorf("查询岗位失败: %w", err)
	}

	validStatuses := map[string]bool{
		model.ApplicationStatusPending:  true,
		model.ApplicationStatusReviewed: true,
		model.ApplicationStatusRejected: true,
		model.ApplicationStatusAccepted: true,
	}

	if !validStatuses[req.NewStatus] {
		return &application.UpdateApplicationStatusResp{
			Success: false,
			Message: "无效的状态值",
		}, nil
	}

	app.Status = req.NewStatus
	if err := s.db.Save(&app).Error; err != nil {
		return nil, fmt.Errorf("更新投递状态失败: %w", err)
	}

	return &application.UpdateApplicationStatusResp{
		Success: true,
		Message: "状态更新成功",
	}, nil
}
