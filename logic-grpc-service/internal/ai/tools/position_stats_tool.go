package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"gorm.io/gorm"
)

type queryPositionStatsInput struct {
	HrId       int64 `json:"hr_id"`
	PositionId int64 `json:"position_id"`
}

func NewQueryPositionStatsTool(db *gorm.DB) tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "query_position_stats",
			Desc: "查询指定岗位的投递统计数据，包括投递人数、各状态分布、学历分布等",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"hr_id":       {Type: "integer", Desc: "HR用户ID"},
				"position_id": {Type: "integer", Desc: "岗位ID"},
			}),
		},
		func(ctx context.Context, input *queryPositionStatsInput) (output string, err error) {
			if input.HrId <= 0 {
				return "", fmt.Errorf("hr_id必须大于0")
			}
			if input.PositionId <= 0 {
				return "", fmt.Errorf("position_id必须大于0")
			}

			var positionExists bool
			err = db.Model(&positionStats{}).
				Where("id = ? AND hr_id = ? AND deleted_at IS NULL", input.PositionId, input.HrId).
				Select("count(*) > 0").Scan(&positionExists).Error
			if err != nil {
				return "", fmt.Errorf("验证岗位权限失败: %w", err)
			}
			if !positionExists {
				return "", fmt.Errorf("岗位不存在或无权限访问")
			}

			var position positionStats
			err = db.Where("id = ? AND deleted_at IS NULL", input.PositionId).First(&position).Error
			if err != nil {
				return "", fmt.Errorf("查询岗位信息失败: %w", err)
			}

			var total int64
			err = db.Model(&applicationStats{}).
				Where("position_id = ? AND deleted_at IS NULL", input.PositionId).
				Count(&total).Error
			if err != nil {
				return "", fmt.Errorf("统计投递总数失败: %w", err)
			}

			var pending, reviewed, accepted, rejected int64

			err = db.Model(&applicationStats{}).
				Where("position_id = ? AND status = ? AND deleted_at IS NULL", input.PositionId, "pending").
				Count(&pending).Error
			if err != nil {
				return "", fmt.Errorf("统计待处理数失败: %w", err)
			}
			err = db.Model(&applicationStats{}).
				Where("position_id = ? AND status = ? AND deleted_at IS NULL", input.PositionId, "reviewed").
				Count(&reviewed).Error
			if err != nil {
				return "", fmt.Errorf("统计已查看数失败: %w", err)
			}
			err = db.Model(&applicationStats{}).
				Where("position_id = ? AND status = ? AND deleted_at IS NULL", input.PositionId, "accepted").
				Count(&accepted).Error
			if err != nil {
				return "", fmt.Errorf("统计通过数失败: %w", err)
			}
			err = db.Model(&applicationStats{}).
				Where("position_id = ? AND status = ? AND deleted_at IS NULL", input.PositionId, "rejected").
				Count(&rejected).Error
			if err != nil {
				return "", fmt.Errorf("统计拒绝数失败: %w", err)
			}

			var educationStats []struct {
				Education string `json:"education"`
				Count     int64  `json:"count"`
			}
			err = db.Model(&userStats{}).
				Joins("JOIN applications ON users.id = applications.candidate_id").
				Where("applications.position_id = ? AND users.education IS NOT NULL AND applications.deleted_at IS NULL AND users.deleted_at IS NULL", input.PositionId).
				Group("users.education").
				Select("users.education, count(*) as count").
				Scan(&educationStats).Error
			if err != nil {
				return "", fmt.Errorf("统计学历分布失败: %w", err)
			}

			result := map[string]interface{}{
				"position_id":        input.PositionId,
				"position_title":     position.Title,
				"location":           position.Location,
				"total_applications": total,
				"status_distribution": map[string]int64{
					"pending":  pending,
					"reviewed": reviewed,
					"accepted": accepted,
					"rejected": rejected,
				},
				"education_distribution": educationStats,
				"created_at":             position.CreatedAt.Format(time.RFC3339),
			}
			data, _ := json.Marshal(result)
			return string(data), nil
		},
	)
}

type positionStats struct {
	ID        int64     `gorm:"column:id"`
	HRID      int64     `gorm:"column:hr_id"`
	Title     string    `gorm:"column:title"`
	Location  *string   `gorm:"column:location"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (positionStats) TableName() string {
	return "positions"
}

type userStats struct {
	ID        int64  `gorm:"column:id"`
	Education string `gorm:"column:education"`
	Skills    string `gorm:"column:skills"`
}

func (userStats) TableName() string {
	return "users"
}
