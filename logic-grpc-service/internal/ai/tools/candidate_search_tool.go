package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"gorm.io/gorm"
)

type queryCandidatesInput struct {
	HrId       int64  `json:"hr_id"`
	Title      string `json:"title"`
	Education  string `json:"education"`
	Skills     string `json:"skills"`
	Experience string `json:"experience"`
	Status     string `json:"status"`
}

func maskPhone(phone *string) *string {
	if phone == nil || *phone == "" {
		return phone
	}
	re := regexp.MustCompile(`^(\d{3})\d{4}(\d{4})$`)
	masked := re.ReplaceAllString(*phone, "$1****$2")
	return &masked
}

func NewQueryCandidatesTool(db *gorm.DB) tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "query_candidates",
			Desc: "按条件筛选候选人，支持按岗位名称、学历、技能、经验、投递状态进行筛选",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"hr_id":      {Type: "integer", Desc: "HR用户ID"},
				"title":      {Type: "string", Desc: "岗位名称（可选，为空则查询所有岗位，支持模糊匹配）"},
				"education":  {Type: "string", Desc: "学历（可选：'高中', '中专', '大专', '本科', '硕士', '博士', '其他'）"},
				"skills":     {Type: "string", Desc: "技能关键词（可选，支持模糊匹配）"},
				"experience": {Type: "string", Desc: "工作经验（可选：1-3年/3-5年/5年以上，支持模糊匹配）"},
				"status":     {Type: "string", Desc: "投递状态（可选：pending/reviewed/accepted/rejected）"},
			}),
		},
		func(ctx context.Context, input *queryCandidatesInput) (output string, err error) {
			if input.HrId <= 0 {
				return "", fmt.Errorf("hr_id必须大于0")
			}

			query := db.Table("applications").
				Select("users.id as user_id, users.username, users.real_name, users.phone, users.education, users.school, users.experience, users.skills, applications.id as application_id, applications.position_id, applications.status, applications.created_at").
				Joins("JOIN users ON applications.candidate_id = users.id").
				Joins("JOIN positions ON applications.position_id = positions.id").
				Where("positions.hr_id = ? AND positions.deleted_at IS NULL AND applications.deleted_at IS NULL AND users.deleted_at IS NULL", input.HrId)

			if input.Title != "" {
				query = query.Where("positions.title LIKE ?", "%"+input.Title+"%")
			}

			if input.Education != "" {
				query = query.Where("users.education = ?", input.Education)
			}

			if input.Skills != "" {
				query = query.Where("users.skills LIKE ?", "%"+input.Skills+"%")
			}

			if input.Experience != "" {
				query = query.Where("users.experience LIKE ?", "%"+input.Experience+"%")
			}

			if input.Status != "" {
				validStatuses := map[string]bool{
					"pending":  true,
					"reviewed": true,
					"accepted": true,
					"rejected": true,
				}
				if !validStatuses[input.Status] {
					return "", fmt.Errorf("无效的状态值，有效值为：pending, reviewed, accepted, rejected")
				}
				query = query.Where("applications.status = ?", input.Status)
			}

			var candidates []struct {
				UserID        int64     `json:"user_id"`
				Username      string    `json:"username"`
				RealName      *string   `json:"real_name"`
				Phone         *string   `json:"phone"`
				Education     *string   `json:"education"`
				School        *string   `json:"school"`
				Experience    *string   `json:"experience"`
				Skills        *string   `json:"skills"`
				ApplicationID int64     `json:"application_id"`
				PositionID    int64     `json:"position_id"`
				Status        string    `json:"status"`
				CreatedAt     time.Time `json:"created_at"`
			}

			err = query.Order("applications.created_at DESC").Limit(50).Scan(&candidates).Error
			if err != nil {
				return "", fmt.Errorf("查询候选人失败: %w", err)
			}

			var positionTitles map[int64]string
			if len(candidates) > 0 {
				positionTitles = make(map[int64]string)
				var positionIDs []int64
				for _, c := range candidates {
					positionIDs = append(positionIDs, c.PositionID)
				}

				var positions []struct {
					ID    int64  `json:"id"`
					Title string `json:"title"`
				}
				err = db.Model(&positionStats{}).
					Where("id IN ?", positionIDs).
					Scan(&positions).Error
				if err != nil {
					return "", fmt.Errorf("查询岗位名称失败: %w", err)
				}

				for _, p := range positions {
					positionTitles[p.ID] = p.Title
				}
			}

			result := make([]map[string]interface{}, 0, len(candidates))
			for _, c := range candidates {
				candidateInfo := map[string]interface{}{
					"user_id":        c.UserID,
					"username":       c.Username,
					"real_name":      c.RealName,
					"phone":          maskPhone(c.Phone),
					"education":      c.Education,
					"school":         c.School,
					"experience":     c.Experience,
					"skills":         c.Skills,
					"application_id": c.ApplicationID,
					"position_id":    c.PositionID,
					"position_title": positionTitles[c.PositionID],
					"status":         c.Status,
					"applied_at":     c.CreatedAt.Format(time.RFC3339),
				}
				result = append(result, candidateInfo)
			}

			response := map[string]interface{}{
				"total_count": len(result),
				"candidates":  result,
				"query_params": map[string]interface{}{
					"hr_id":      input.HrId,
					"title":      input.Title,
					"education":  input.Education,
					"skills":     input.Skills,
					"experience": input.Experience,
					"status":     input.Status,
				},
			}

			data, _ := json.Marshal(response)
			return string(data), nil
		},
	)
}
