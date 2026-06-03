package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"gorm.io/gorm"
)

type applicationStats struct {
	ID          int64          `gorm:"column:id"`
	PositionID  int64          `gorm:"column:position_id"`
	CandidateID int64          `gorm:"column:candidate_id"`
	Status      string         `gorm:"column:status"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (applicationStats) TableName() string {
	return "applications"
}

type queryApplicationStatsInput struct {
	HrId int64 `json:"hr_id"`
}

func NewQueryApplicationStatsTool(db *gorm.DB) tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "query_application_stats",
			Desc: "查询HR用户的投递统计数据，包括总投递数、待处理数、已查看数、通过数、拒绝数",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"hr_id": {Type: "integer", Desc: "HR用户ID"},
			}),
		},
		func(ctx context.Context, input *queryApplicationStatsInput) (output string, err error) {
			if input.HrId <= 0 {
				return "", fmt.Errorf("hr_id必须大于0")
			}

			var total int64
			err = db.Model(&applicationStats{}).
				Joins("JOIN positions ON applications.position_id = positions.id").
				Where("positions.hr_id = ? AND applications.deleted_at IS NULL", input.HrId).
				Count(&total).Error
			if err != nil {
				return "", fmt.Errorf("统计总投递数失败: %w", err)
			}

			var pending int64
			err = db.Model(&applicationStats{}).
				Joins("JOIN positions ON applications.position_id = positions.id").
				Where("positions.hr_id = ? AND applications.status = ? AND applications.deleted_at IS NULL", input.HrId, "pending").
				Count(&pending).Error
			if err != nil {
				return "", fmt.Errorf("统计待处理数失败: %w", err)
			}

			var reviewed int64
			err = db.Model(&applicationStats{}).
				Joins("JOIN positions ON applications.position_id = positions.id").
				Where("positions.hr_id = ? AND applications.status = ? AND applications.deleted_at IS NULL", input.HrId, "reviewed").
				Count(&reviewed).Error
			if err != nil {
				return "", fmt.Errorf("统计已查看数失败: %w", err)
			}

			var accepted int64
			err = db.Model(&applicationStats{}).
				Joins("JOIN positions ON applications.position_id = positions.id").
				Where("positions.hr_id = ? AND applications.status = ? AND applications.deleted_at IS NULL", input.HrId, "accepted").
				Count(&accepted).Error
			if err != nil {
				return "", fmt.Errorf("统计通过数失败: %w", err)
			}

			var rejected int64
			err = db.Model(&applicationStats{}).
				Joins("JOIN positions ON applications.position_id = positions.id").
				Where("positions.hr_id = ? AND applications.status = ? AND applications.deleted_at IS NULL", input.HrId, "rejected").
				Count(&rejected).Error
			if err != nil {
				return "", fmt.Errorf("统计拒绝数失败: %w", err)
			}

			result := map[string]interface{}{
				"hr_id":    input.HrId,
				"total":    total,
				"pending":  pending,
				"reviewed": reviewed,
				"accepted": accepted,
				"rejected": rejected,
			}
			data, _ := json.Marshal(result)
			return string(data), nil
		},
	)
}
