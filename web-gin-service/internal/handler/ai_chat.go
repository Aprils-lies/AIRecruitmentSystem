package handler

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"

	"web-gin-service/internal/grpc"
	"web-gin-service/internal/middleware"
	"web-gin-service/pkg"
	"web-gin-service/proto/gen/ai_chat"
)

type AIChatHandler struct{}

func NewAIChatHandler() *AIChatHandler {
	return &AIChatHandler{}
}

func (h *AIChatHandler) Chat(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	var req struct {
		Question  string `json:"question" binding:"required"`
		SessionID string `json:"session_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ParamError(c, "参数错误："+err.Error())
		return
	}

	grpcReq := &ai_chat.ChatReq{
		HrId:      hrID,
		Question:  req.Question,
		SessionId: req.SessionID,
	}

	resp, err := grpc.GetAIChatClient().Chat(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "AI对话失败："+err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"answer":     pkg.SafeString(resp.Answer),
		"message_id": resp.MessageId,
		"session_id": pkg.SafeString(resp.SessionId),
	})
}

func (h *AIChatHandler) ChatStream(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	var req struct {
		Question  string `json:"question" binding:"required"`
		SessionID string `json:"session_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ParamError(c, "参数错误："+err.Error())
		return
	}

	grpcReq := &ai_chat.ChatReq{
		HrId:      hrID,
		Question:  req.Question,
		SessionId: req.SessionID,
	}

	stream, err := grpc.GetAIChatClient().ChatStream(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "AI对话失败："+err.Error())
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}

		data := map[string]interface{}{
			"chunk":      pkg.SafeString(resp.Chunk),
			"done":       resp.Done,
			"message_id": resp.MessageId,
			"session_id": pkg.SafeString(resp.SessionId),
		}

		jsonData, _ := json.Marshal(data)
		c.Writer.WriteString("data: " + string(jsonData) + "\n\n")
		c.Writer.Flush()

		if resp.Done {
			break
		}
	}
}

func (h *AIChatHandler) GetHistory(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	sessionID := c.Query("session_id")
	if sessionID == "" {
		pkg.ParamError(c, "session_id不能为空")
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)

	grpcReq := &ai_chat.GetHistoryReq{
		HrId:      hrID,
		SessionId: sessionID,
		Limit:     int32(limit),
		Offset:    int32(offset),
	}

	resp, err := grpc.GetAIChatClient().GetHistory(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取对话历史失败："+err.Error())
		return
	}

	items := make([]map[string]interface{}, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, map[string]interface{}{
			"id":         item.Id,
			"session_id": pkg.SafeString(item.SessionId),
			"role":       pkg.SafeString(item.Role),
			"content":    pkg.SafeString(item.Content),
			"created_at": pkg.SafeString(item.CreatedAt),
		})
	}

	pkg.Success(c, gin.H{
		"items": items,
		"total": resp.Total,
	})
}

func (h *AIChatHandler) ListSessions(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	grpcReq := &ai_chat.ListSessionsReq{
		HrId: hrID,
	}

	resp, err := grpc.GetAIChatClient().ListSessions(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取会话列表失败："+err.Error())
		return
	}

	sessions := make([]map[string]interface{}, 0, len(resp.Sessions))
	for _, session := range resp.Sessions {
		sessions = append(sessions, map[string]interface{}{
			"session_id":   pkg.SafeString(session.SessionId),
			"last_message": pkg.SafeString(session.LastMessage),
			"created_at":   pkg.SafeString(session.CreatedAt),
		})
	}

	pkg.Success(c, gin.H{
		"sessions": sessions,
	})
}

func (h *AIChatHandler) GetStats(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	grpcReq := &ai_chat.GetStatsReq{
		HrId: hrID,
	}

	resp, err := grpc.GetAIChatClient().GetStats(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取统计数据失败："+err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"summary": pkg.SafeString(resp.Summary),
	})
}

func (h *AIChatHandler) DeleteSession(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		pkg.ParamError(c, "session_id不能为空")
		return
	}

	grpcReq := &ai_chat.DeleteSessionReq{
		HrId:      hrID,
		SessionId: sessionID,
	}

	resp, err := grpc.GetAIChatClient().DeleteSession(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "删除会话失败："+err.Error())
		return
	}

	if !resp.Success {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, nil)
}

func (h *AIChatHandler) DeleteMessage(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	messageIDStr := c.Param("message_id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		pkg.ParamError(c, "无效的消息ID")
		return
	}

	grpcReq := &ai_chat.DeleteMessageReq{
		HrId:      hrID,
		MessageId: messageID,
	}

	resp, err := grpc.GetAIChatClient().DeleteMessage(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "删除消息失败："+err.Error())
		return
	}

	if !resp.Success {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, nil)
}
