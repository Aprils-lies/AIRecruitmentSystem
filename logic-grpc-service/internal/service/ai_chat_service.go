package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"

	"logic-grpc-service/internal/ai"
	"logic-grpc-service/internal/model"
	"logic-grpc-service/proto/gen/ai_chat"
)

type AIChatService struct {
	ai_chat.UnimplementedAIChatServiceServer
	db *gorm.DB
}

func NewAIChatService(db *gorm.DB) *AIChatService {
	return &AIChatService{db: db}
}

func (s *AIChatService) Chat(ctx context.Context, req *ai_chat.ChatReq) (*ai_chat.ChatResp, error) {
	var history []model.ChatHistory
	if req.SessionId != "" {
		if err := s.db.Where("hr_id = ? AND session_id = ?", req.HrId, req.SessionId).Order("created_at ASC, id ASC").Find(&history).Error; err != nil {
			log.Printf("查询对话历史失败: hr_id=%d, session_id=%s, err=%v", req.HrId, req.SessionId, err)
			return nil, fmt.Errorf("查询对话历史失败")
		}
	}

	sessionID := req.SessionId
	if sessionID == "" {
		sessionID = generateSessionID()
	}

	userMsg := model.ChatHistory{
		HRID:      req.HrId,
		SessionID: sessionID,
		Role:      model.ChatRoleUser,
		Content:   req.Question,
	}
	if err := s.db.Create(&userMsg).Error; err != nil {
		log.Printf("保存用户消息失败: hr_id=%d, err=%v", req.HrId, err)
		return nil, fmt.Errorf("保存消息失败")
	}

	var messages []ai.ChatMessage
	for _, h := range history {
		messages = append(messages, ai.ChatMessage{
			Role:    h.Role,
			Content: h.Content,
		})
	}

	messages = append(messages, ai.ChatMessage{
		Role:    model.ChatRoleUser,
		Content: req.Question,
		HrId:    req.HrId,
	})

	answer, err := ai.Chat(ctx, messages)
	if err != nil {
		return nil, err
	}

	assistantMsg := model.ChatHistory{
		HRID:      req.HrId,
		SessionID: sessionID,
		Role:      model.ChatRoleAssistant,
		Content:   answer,
	}
	if err := s.db.Create(&assistantMsg).Error; err != nil {
		log.Printf("保存AI回复失败: hr_id=%d, err=%v", req.HrId, err)
		return nil, fmt.Errorf("保存消息失败")
	}

	return &ai_chat.ChatResp{
		Answer:    answer,
		MessageId: assistantMsg.ID,
		SessionId: sessionID,
	}, nil
}

func (s *AIChatService) ChatStream(req *ai_chat.ChatReq, stream ai_chat.AIChatService_ChatStreamServer) error {
	ctx := stream.Context()

	sessionID := req.SessionId
	if sessionID == "" {
		sessionID = generateSessionID()
	}

	var history []model.ChatHistory
	if req.SessionId != "" {
		if err := s.db.Where("hr_id = ? AND session_id = ?", req.HrId, req.SessionId).Order("created_at ASC, id ASC").Find(&history).Error; err != nil {
			log.Printf("查询对话历史失败: hr_id=%d, session_id=%s, err=%v", req.HrId, req.SessionId, err)
			return fmt.Errorf("查询对话历史失败")
		}
	}

	userMsg := model.ChatHistory{
		HRID:      req.HrId,
		SessionID: sessionID,
		Role:      model.ChatRoleUser,
		Content:   req.Question,
	}
	if err := s.db.Create(&userMsg).Error; err != nil {
		log.Printf("保存用户消息失败: hr_id=%d, err=%v", req.HrId, err)
		return fmt.Errorf("保存消息失败")
	}

	var messages []ai.ChatMessage
	for _, h := range history {
		messages = append(messages, ai.ChatMessage{
			Role:    h.Role,
			Content: h.Content,
		})
	}

	messages = append(messages, ai.ChatMessage{
		Role:    model.ChatRoleUser,
		Content: req.Question,
		HrId:    req.HrId,
	})

	streamCh, err := ai.ChatStream(ctx, messages)
	if err != nil {
		return err
	}

	var fullAnswer strings.Builder
	var messageID int64

	for chunk := range streamCh {
		fullAnswer.WriteString(chunk)
		if err := stream.Send(&ai_chat.ChatStreamResp{
			Chunk:     chunk,
			Done:      false,
			MessageId: 0,
			SessionId: sessionID,
		}); err != nil {
			log.Printf("发送流式响应失败: hr_id=%d, err=%v", req.HrId, err)
			return fmt.Errorf("发送消息失败")
		}
	}

	assistantMsg := model.ChatHistory{
		HRID:      req.HrId,
		SessionID: sessionID,
		Role:      model.ChatRoleAssistant,
		Content:   fullAnswer.String(),
	}
	if err := s.db.Create(&assistantMsg).Error; err != nil {
		log.Printf("保存AI回复失败: hr_id=%d, err=%v", req.HrId, err)
		return fmt.Errorf("保存消息失败")
	}
	messageID = assistantMsg.ID

	if err := stream.Send(&ai_chat.ChatStreamResp{
		Chunk:     "",
		Done:      true,
		MessageId: messageID,
		SessionId: sessionID,
	}); err != nil {
		log.Printf("发送完成信号失败: hr_id=%d, err=%v", req.HrId, err)
		return fmt.Errorf("发送消息失败")
	}

	return nil
}

func (s *AIChatService) GetHistory(ctx context.Context, req *ai_chat.GetHistoryReq) (*ai_chat.GetHistoryResp, error) {
	query := s.db.Model(&model.ChatHistory{}).Where("hr_id = ? AND session_id = ?", req.HrId, req.SessionId)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Printf("统计消息数量失败: hr_id=%d, session_id=%s, err=%v", req.HrId, req.SessionId, err)
		return nil, fmt.Errorf("查询失败")
	}

	var history []model.ChatHistory
	if err := query.Order("created_at DESC, id DESC").Offset(int(req.Offset)).Limit(int(req.Limit)).Find(&history).Error; err != nil {
		log.Printf("查询对话历史失败: hr_id=%d, session_id=%s, err=%v", req.HrId, req.SessionId, err)
		return nil, fmt.Errorf("查询失败")
	}

	result := make([]*ai_chat.HistoryItem, 0, len(history))
	for i := len(history) - 1; i >= 0; i-- {
		h := history[i]
		result = append(result, &ai_chat.HistoryItem{
			Id:        h.ID,
			SessionId: h.SessionID,
			Role:      h.Role,
			Content:   h.Content,
			CreatedAt: h.CreatedAt.Format(time.RFC3339),
		})
	}

	return &ai_chat.GetHistoryResp{
		Items: result,
		Total: int32(total),
	}, nil
}

func (s *AIChatService) ListSessions(ctx context.Context, req *ai_chat.ListSessionsReq) (*ai_chat.ListSessionsResp, error) {
	var sessions []struct {
		SessionID   string    `gorm:"column:session_id"`
		LastMessage string    `gorm:"column:content"`
		CreatedAt   time.Time `gorm:"column:created_at"`
	}

	sql := `
		SELECT session_id, content, created_at
		FROM chat_history
		WHERE hr_id = ?
		AND deleted_at IS NULL
		AND id IN (
			SELECT MAX(id) FROM chat_history WHERE hr_id = ? AND deleted_at IS NULL GROUP BY session_id
		)
		ORDER BY created_at DESC, id DESC
	`
	if err := s.db.Raw(sql, req.HrId, req.HrId).Scan(&sessions).Error; err != nil {
		log.Printf("查询会话列表失败: hr_id=%d, err=%v", req.HrId, err)
		return nil, fmt.Errorf("查询失败")
	}

	result := make([]*ai_chat.SessionInfo, 0, len(sessions))
	for _, s := range sessions {
		result = append(result, &ai_chat.SessionInfo{
			SessionId:   s.SessionID,
			LastMessage: s.LastMessage,
			CreatedAt:   s.CreatedAt.Format(time.RFC3339),
		})
	}

	return &ai_chat.ListSessionsResp{
		Sessions: result,
	}, nil
}

func (s *AIChatService) GetStats(ctx context.Context, req *ai_chat.GetStatsReq) (*ai_chat.GetStatsResp, error) {
	var totalPositions int64
	if err := s.db.Model(&model.Position{}).Where("hr_id = ?", req.HrId).Count(&totalPositions).Error; err != nil {
		log.Printf("统计岗位数量失败: hr_id=%d, err=%v", req.HrId, err)
		return nil, fmt.Errorf("统计失败")
	}

	var totalApplications int64
	if err := s.db.Model(&model.Application{}).
		Joins("JOIN positions ON applications.position_id = positions.id").
		Where("positions.hr_id = ?", req.HrId).
		Count(&totalApplications).Error; err != nil {
		log.Printf("统计投递数量失败: hr_id=%d, err=%v", req.HrId, err)
		return nil, fmt.Errorf("统计失败")
	}

	var pendingCount int64
	if err := s.db.Model(&model.Application{}).
		Joins("JOIN positions ON applications.position_id = positions.id").
		Where("positions.hr_id = ? AND applications.status = ?", req.HrId, model.ApplicationStatusPending).
		Count(&pendingCount).Error; err != nil {
		log.Printf("统计待处理数量失败: hr_id=%d, err=%v", req.HrId, err)
		return nil, fmt.Errorf("统计失败")
	}

	var reviewedCount int64
	if err := s.db.Model(&model.Application{}).
		Joins("JOIN positions ON applications.position_id = positions.id").
		Where("positions.hr_id = ? AND applications.status = ?", req.HrId, model.ApplicationStatusReviewed).
		Count(&reviewedCount).Error; err != nil {
		log.Printf("统计已查看数量失败: hr_id=%d, err=%v", req.HrId, err)
		return nil, fmt.Errorf("统计失败")
	}

	var acceptedCount int64
	if err := s.db.Model(&model.Application{}).
		Joins("JOIN positions ON applications.position_id = positions.id").
		Where("positions.hr_id = ? AND applications.status = ?", req.HrId, model.ApplicationStatusAccepted).
		Count(&acceptedCount).Error; err != nil {
		log.Printf("统计通过数量失败: hr_id=%d, err=%v", req.HrId, err)
		return nil, fmt.Errorf("统计失败")
	}

	summary := fmt.Sprintf("您共有 %d 个岗位，收到 %d 份投递。其中待处理 %d 份，已查看 %d 份，通过 %d 份。",
		totalPositions, totalApplications, pendingCount, reviewedCount, acceptedCount)

	return &ai_chat.GetStatsResp{
		Summary: summary,
	}, nil
}

func (s *AIChatService) DeleteSession(ctx context.Context, req *ai_chat.DeleteSessionReq) (*ai_chat.DeleteSessionResp, error) {
	if err := s.db.Where("hr_id = ? AND session_id = ?", req.HrId, req.SessionId).Delete(&model.ChatHistory{}).Error; err != nil {
		log.Printf("删除会话失败: hr_id=%d, session_id=%s, err=%v", req.HrId, req.SessionId, err)
		return nil, fmt.Errorf("删除失败")
	}

	return &ai_chat.DeleteSessionResp{
		Success: true,
		Message: "删除成功",
	}, nil
}

func (s *AIChatService) DeleteMessage(ctx context.Context, req *ai_chat.DeleteMessageReq) (*ai_chat.DeleteMessageResp, error) {
	if err := s.db.Where("id = ? AND hr_id = ?", req.MessageId, req.HrId).Delete(&model.ChatHistory{}).Error; err != nil {
		log.Printf("删除消息失败: hr_id=%d, message_id=%d, err=%v", req.HrId, req.MessageId, err)
		return nil, fmt.Errorf("删除失败")
	}

	return &ai_chat.DeleteMessageResp{
		Success: true,
		Message: "删除成功",
	}, nil
}

func generateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
