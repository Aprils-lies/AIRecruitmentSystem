package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"

	"web-gin-service/internal/grpc"
	"web-gin-service/internal/middleware"
	"web-gin-service/pkg"
	"web-gin-service/proto/gen/application"
)

type ApplicationHandler struct{}

func NewApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{}
}

func (h *ApplicationHandler) Apply(c *gin.Context) {
	candidateID := middleware.GetUserID(c)
	if candidateID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	positionIDStr := c.Param("id")
	positionID, err := strconv.ParseInt(positionIDStr, 10, 64)
	if err != nil {
		pkg.ParamError(c, "无效的岗位ID")
		return
	}

	var req struct {
		ResumeID int64 `json:"resume_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ParamError(c, "参数错误："+err.Error())
		return
	}

	grpcReq := &application.ApplyReq{
		PositionId:  positionID,
		CandidateId: candidateID,
		ResumeId:   req.ResumeID,
	}

	resp, err := grpc.GetApplicationClient().Apply(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "投递失败："+err.Error())
		return
	}

	if resp.ApplicationId == 0 {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, gin.H{
		"application_id": resp.ApplicationId,
	})
}

func (h *ApplicationHandler) ListCandidates(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	positionIDStr := c.Param("id")
	positionID, err := strconv.ParseInt(positionIDStr, 10, 64)
	if err != nil {
		pkg.ParamError(c, "无效的岗位ID")
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)

	grpcReq := &application.ListCandidatesReq{
		PositionId: positionID,
		HrId:       hrID,
		Page:       int32(page),
		PageSize:   int32(pageSize),
	}

	resp, err := grpc.GetApplicationClient().ListCandidates(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取候选人列表失败："+err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"candidates": resp.Candidates,
		"total":     resp.Total,
		"page":      resp.Page,
	})
}

func (h *ApplicationHandler) GetCandidateDetail(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	candidateIDStr := c.Param("id")
	candidateID, err := strconv.ParseInt(candidateIDStr, 10, 64)
	if err != nil {
		pkg.ParamError(c, "无效的候选人ID")
		return
	}

	applicationIDStr := c.Query("application_id")
	applicationID, err := strconv.ParseInt(applicationIDStr, 10, 64)
	if err != nil {
		pkg.ParamError(c, "无效的投递记录ID")
		return
	}

	grpcReq := &application.GetCandidateDetailReq{
		CandidateId:   candidateID,
		HrId:          hrID,
		ApplicationId: applicationID,
	}

	resp, err := grpc.GetApplicationClient().GetCandidateDetail(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取候选人详情失败："+err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"user_id":     resp.UserId,
		"username":    pkg.SafeString(resp.Username),
		"real_name":   pkg.SafeString(resp.RealName),
		"phone":       pkg.SafeString(resp.Phone),
		"education":   pkg.SafeString(resp.Education),
		"school":      pkg.SafeString(resp.School),
		"experience":  pkg.SafeString(resp.Experience),
		"skills":      pkg.SafeString(resp.Skills),
		"resume_id":   resp.ResumeId,
		"resume_name": pkg.SafeString(resp.ResumeName),
		"resume_type": pkg.SafeString(resp.ResumeType),
		"applied_at":  pkg.SafeString(resp.AppliedAt),
	})
}

func (h *ApplicationHandler) ListMyApplications(c *gin.Context) {
	candidateID := middleware.GetUserID(c)
	if candidateID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)

	grpcReq := &application.ListMyApplicationsReq{
		CandidateId: candidateID,
		Page:        int32(page),
		PageSize:    int32(pageSize),
	}

	resp, err := grpc.GetApplicationClient().ListMyApplications(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取投递记录失败："+err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"applications": resp.Applications,
		"total":       resp.Total,
		"page":        resp.Page,
	})
}

func (h *ApplicationHandler) WithdrawApplication(c *gin.Context) {
	candidateID := middleware.GetUserID(c)
	if candidateID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	applicationIDStr := c.Param("id")
	applicationID, err := strconv.ParseInt(applicationIDStr, 10, 64)
	if err != nil {
		pkg.ParamError(c, "无效的投递记录ID")
		return
	}

	grpcReq := &application.WithdrawApplicationReq{
		ApplicationId: applicationID,
		CandidateId:   candidateID,
	}

	resp, err := grpc.GetApplicationClient().WithdrawApplication(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "撤回投递失败："+err.Error())
		return
	}

	if !resp.Success {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, nil)
}

func (h *ApplicationHandler) UpdateApplicationStatus(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	applicationIDStr := c.Param("id")
	applicationID, err := strconv.ParseInt(applicationIDStr, 10, 64)
	if err != nil {
		pkg.ParamError(c, "无效的投递记录ID")
		return
	}

	var req struct {
		NewStatus string `json:"new_status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ParamError(c, "参数错误："+err.Error())
		return
	}

	grpcReq := &application.UpdateApplicationStatusReq{
		ApplicationId: applicationID,
		HrId:          hrID,
		NewStatus:     req.NewStatus,
	}

	resp, err := grpc.GetApplicationClient().UpdateApplicationStatus(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "更新投递状态失败："+err.Error())
		return
	}

	if !resp.Success {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, nil)
}
