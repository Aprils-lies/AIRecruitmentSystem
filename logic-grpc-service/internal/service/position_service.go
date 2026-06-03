package service

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"logic-grpc-service/internal/model"
	"logic-grpc-service/proto/gen/position"
)

type PositionService struct {
	position.UnimplementedPositionServiceServer
	db *gorm.DB
}

func NewPositionService(db *gorm.DB) *PositionService {
	return &PositionService{db: db}
}

func (s *PositionService) CreatePosition(ctx context.Context, req *position.CreatePositionReq) (*position.CreatePositionResp, error) {
	positionModel := model.Position{
		HRID:        req.HrId,
		Title:       req.Title,
		Description: req.Description,
		Requirements: req.Requirements,
		SalaryMin:   &req.SalaryMin,
		SalaryMax:   &req.SalaryMax,
		Location:    &req.Location,
		Status:      model.PositionStatusPublished,
	}

	if err := s.db.Create(&positionModel).Error; err != nil {
		return nil, fmt.Errorf("创建岗位失败: %w", err)
	}

	return &position.CreatePositionResp{
		PositionId: positionModel.ID,
		Message:    "创建成功",
	}, nil
}

func (s *PositionService) UpdatePosition(ctx context.Context, req *position.UpdatePositionReq) (*position.UpdatePositionResp, error) {
	var pos model.Position
	if err := s.db.Where("id = ? AND hr_id = ?", req.PositionId, req.HrId).First(&pos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &position.UpdatePositionResp{
				Success: false,
				Message: "岗位不存在或无权限",
			}, nil
		}
		return nil, fmt.Errorf("查询岗位失败: %w", err)
	}

	updateData := map[string]interface{}{}
	if req.Title != "" {
		updateData["title"] = req.Title
	}
	if req.Description != "" {
		updateData["description"] = req.Description
	}
	if req.Requirements != "" {
		updateData["requirements"] = req.Requirements
	}
	updateData["salary_min"] = req.SalaryMin
	updateData["salary_max"] = req.SalaryMax
	if req.Location != "" {
		updateData["location"] = req.Location
	}

	if err := s.db.Model(&pos).Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("更新岗位失败: %w", err)
	}

	return &position.UpdatePositionResp{
		Success: true,
		Message: "更新成功",
	}, nil
}

func (s *PositionService) OfflinePosition(ctx context.Context, req *position.OfflinePositionReq) (*position.OfflinePositionResp, error) {
	var pos model.Position
	if err := s.db.Where("id = ? AND hr_id = ?", req.PositionId, req.HrId).First(&pos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &position.OfflinePositionResp{
				Success: false,
				Message: "岗位不存在或无权限",
			}, nil
		}
		return nil, fmt.Errorf("查询岗位失败: %w", err)
	}

	pos.Status = model.PositionStatusOffline
	if err := s.db.Save(&pos).Error; err != nil {
		return nil, fmt.Errorf("下架岗位失败: %w", err)
	}

	return &position.OfflinePositionResp{
		Success: true,
		Message: "已下架",
	}, nil
}

func (s *PositionService) ListPositions(ctx context.Context, req *position.ListPositionsReq) (*position.ListPositionsResp, error) {
	query := s.db.Model(&model.Position{}).Where("status = ?", model.PositionStatusPublished)

	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Location != "" {
		query = query.Where("location = ?", req.Location)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("统计岗位数量失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	var positions []model.Position
	if err := query.Order("created_at DESC").Offset(int(offset)).Limit(int(req.PageSize)).Find(&positions).Error; err != nil {
		return nil, fmt.Errorf("查询岗位列表失败: %w", err)
	}

	result := make([]*position.PositionInfo, 0, len(positions))
	for _, pos := range positions {
		result = append(result, &position.PositionInfo{
			Id:           pos.ID,
			HrId:         pos.HRID,
			Title:        pos.Title,
			Description:  pos.Description,
			Requirements: pos.Requirements,
			SalaryMin:    getInt32Ptr(pos.SalaryMin),
			SalaryMax:    getInt32Ptr(pos.SalaryMax),
			Location:     getStringPtr(pos.Location),
			Status:       pos.Status,
			CreatedAt:    pos.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    pos.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &position.ListPositionsResp{
		Positions: result,
		Total:     int32(total),
		Page:      req.Page,
	}, nil
}

func (s *PositionService) GetPosition(ctx context.Context, req *position.GetPositionReq) (*position.GetPositionResp, error) {
	var pos model.Position
	if err := s.db.Where("id = ?", req.PositionId).First(&pos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &position.GetPositionResp{}, nil
		}
		return nil, fmt.Errorf("查询岗位失败: %w", err)
	}

	return &position.GetPositionResp{
		Position: &position.PositionInfo{
			Id:           pos.ID,
			HrId:         pos.HRID,
			Title:        pos.Title,
			Description:  pos.Description,
			Requirements: pos.Requirements,
			SalaryMin:    getInt32Ptr(pos.SalaryMin),
			SalaryMax:    getInt32Ptr(pos.SalaryMax),
			Location:     getStringPtr(pos.Location),
			Status:       pos.Status,
			CreatedAt:    pos.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    pos.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *PositionService) OnlinePosition(ctx context.Context, req *position.OnlinePositionReq) (*position.OnlinePositionResp, error) {
	var pos model.Position
	if err := s.db.Where("id = ? AND hr_id = ?", req.PositionId, req.HrId).First(&pos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &position.OnlinePositionResp{
				Success: false,
				Message: "岗位不存在或无权限",
			}, nil
		}
		return nil, fmt.Errorf("查询岗位失败: %w", err)
	}

	pos.Status = model.PositionStatusPublished
	if err := s.db.Save(&pos).Error; err != nil {
		return nil, fmt.Errorf("上架岗位失败: %w", err)
	}

	return &position.OnlinePositionResp{
		Success: true,
		Message: "已上架",
	}, nil
}

func (s *PositionService) ListMyPositions(ctx context.Context, req *position.ListMyPositionsReq) (*position.ListMyPositionsResp, error) {
	query := s.db.Model(&model.Position{}).Where("hr_id = ?", req.HrId)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("统计岗位数量失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	var positions []model.Position
	if err := query.Order("created_at DESC").Offset(int(offset)).Limit(int(req.PageSize)).Find(&positions).Error; err != nil {
		return nil, fmt.Errorf("查询岗位列表失败: %w", err)
	}

	result := make([]*position.PositionInfo, 0, len(positions))
	for _, pos := range positions {
		result = append(result, &position.PositionInfo{
			Id:           pos.ID,
			HrId:         pos.HRID,
			Title:        pos.Title,
			Description:  pos.Description,
			Requirements: pos.Requirements,
			SalaryMin:    getInt32Ptr(pos.SalaryMin),
			SalaryMax:    getInt32Ptr(pos.SalaryMax),
			Location:     getStringPtr(pos.Location),
			Status:       pos.Status,
			CreatedAt:    pos.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    pos.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &position.ListMyPositionsResp{
		Positions: result,
		Total:     int32(total),
		Page:      req.Page,
	}, nil
}

func getInt32Ptr(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}