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

type queryPositionHotnessInput struct {
	HrId  int64 `json:"hr_id"`
	Limit int   `json:"limit"`
}

func NewQueryPositionHotnessTool(db *gorm.DB) tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "query_position_hotness",
			Desc: "查询HR用户所有岗位的投递热度排行，按投递数从高到低排序",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"hr_id": {Type: "integer", Desc: "HR用户ID"},
				"limit": {Type: "integer", Desc: "返回数量（可选，默认10）"},
			}),
		},
		func(ctx context.Context, input *queryPositionHotnessInput) (output string, err error) {
			if input.HrId <= 0 {
				return "", fmt.Errorf("hr_id必须大于0")
			}

			limit := input.Limit
			if limit <= 0 {
				limit = 10
			}

			var hotnessList []struct {
				PositionID       int64   `json:"position_id"`
				PositionTitle    string  `json:"position_title"`
				Location         *string `json:"location"`
				ApplicationCount int64   `json:"application_count"`
				Status           string  `json:"status"`
			}

			sql := `
				SELECT 
					p.id as position_id,
					p.title as position_title,
					p.location,
					p.status,
					COALESCE(COUNT(a.id), 0) as application_count
				FROM positions p
				LEFT JOIN applications a ON p.id = a.position_id AND a.deleted_at IS NULL
				WHERE p.hr_id = ? AND p.deleted_at IS NULL
				GROUP BY p.id, p.title, p.location, p.status
				ORDER BY application_count DESC
				LIMIT ?
			`

			err = db.Raw(sql, input.HrId, limit).Scan(&hotnessList).Error
			if err != nil {
				return "", fmt.Errorf("查询岗位热度失败: %w", err)
			}

			var totalPositions int64
			err = db.Model(&positionStats{}).
				Where("hr_id = ?", input.HrId).
				Count(&totalPositions).Error
			if err != nil {
				return "", fmt.Errorf("统计岗位总数失败: %w", err)
			}

			var totalApplications int64
			err = db.Model(&applicationStats{}).
				Joins("JOIN positions ON applications.position_id = positions.id").
				Where("positions.hr_id = ?", input.HrId).
				Count(&totalApplications).Error
			if err != nil {
				return "", fmt.Errorf("统计总投递数失败: %w", err)
			}

			result := map[string]interface{}{
				"hr_id":              input.HrId,
				"total_positions":    totalPositions,
				"total_applications": totalApplications,
				"hotness_list":       hotnessList,
			}
			data, _ := json.Marshal(result)
			return string(data), nil
		},
	)
}
