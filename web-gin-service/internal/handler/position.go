package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"

	"web-gin-service/internal/grpc"
	"web-gin-service/internal/middleware"
	"web-gin-service/pkg"
	"web-gin-service/proto/gen/position"
)

type PositionHandler struct{}

func NewPositionHandler() *PositionHandler {
	return &PositionHandler{}
}

func (h *PositionHandler) CreatePosition(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	var req struct {
		Title        string `json:"title" binding:"required"`
		Description  string `json:"description" binding:"required"`
		Requirements string `json:"requirements"`
		SalaryMin    int32  `json:"salary_min"`
		SalaryMax    int32  `json:"salary_max"`
		Location     string `json:"location"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ParamError(c, "参数错误："+err.Error())
		return
	}

	grpcReq := &position.CreatePositionReq{
		HrId:        hrID,
		Title:       req.Title,
		Description: req.Description,
		Requirements: req.Requirements,
		SalaryMin:   req.SalaryMin,
		SalaryMax:   req.SalaryMax,
		Location:    req.Location,
	}

	resp, err := grpc.GetPositionClient().CreatePosition(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "创建岗位失败："+err.Error())
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, gin.H{
		"position_id": resp.PositionId,
	})
}

func (h *PositionHandler) UpdatePosition(c *gin.Context) {
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

	var req struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		Requirements string `json:"requirements"`
		SalaryMin    int32  `json:"salary_min"`
		SalaryMax    int32  `json:"salary_max"`
		Location     string `json:"location"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ParamError(c, "参数错误："+err.Error())
		return
	}

	grpcReq := &position.UpdatePositionReq{
		PositionId:  positionID,
		HrId:        hrID,
		Title:       req.Title,
		Description: req.Description,
		Requirements: req.Requirements,
		SalaryMin:   req.SalaryMin,
		SalaryMax:   req.SalaryMax,
		Location:    req.Location,
	}

	resp, err := grpc.GetPositionClient().UpdatePosition(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "更新岗位失败："+err.Error())
		return
	}

	if !resp.Success {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, nil)
}

func (h *PositionHandler) OfflinePosition(c *gin.Context) {
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

	grpcReq := &position.OfflinePositionReq{
		PositionId: positionID,
		HrId:       hrID,
	}

	resp, err := grpc.GetPositionClient().OfflinePosition(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "下架岗位失败："+err.Error())
		return
	}

	if !resp.Success {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, nil)
}

func (h *PositionHandler) OnlinePosition(c *gin.Context) {
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

	grpcReq := &position.OnlinePositionReq{
		PositionId: positionID,
		HrId:       hrID,
	}

	resp, err := grpc.GetPositionClient().OnlinePosition(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "上架岗位失败："+err.Error())
		return
	}

	if !resp.Success {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, nil)
}

func (h *PositionHandler) ListPositions(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)
	keyword := c.Query("keyword")
	location := c.Query("location")

	grpcReq := &position.ListPositionsReq{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Keyword:  keyword,
		Location: location,
	}

	resp, err := grpc.GetPositionClient().ListPositions(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取岗位列表失败："+err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"positions": resp.Positions,
		"total":     resp.Total,
		"page":      resp.Page,
	})
}

func (h *PositionHandler) GetPosition(c *gin.Context) {
	positionIDStr := c.Param("id")
	positionID, err := strconv.ParseInt(positionIDStr, 10, 64)
	if err != nil {
		pkg.ParamError(c, "无效的岗位ID")
		return
	}

	grpcReq := &position.GetPositionReq{
		PositionId: positionID,
	}

	resp, err := grpc.GetPositionClient().GetPosition(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取岗位详情失败："+err.Error())
		return
	}

	if resp.Position == nil {
		pkg.NotFound(c, "岗位不存在")
		return
	}

	pkg.Success(c, resp.Position)
}

func (h *PositionHandler) ListMyPositions(c *gin.Context) {
	hrID := middleware.GetUserID(c)
	if hrID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)

	grpcReq := &position.ListMyPositionsReq{
		HrId:     hrID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	resp, err := grpc.GetPositionClient().ListMyPositions(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取我的岗位列表失败："+err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"positions": resp.Positions,
		"total":     resp.Total,
		"page":      resp.Page,
	})
}
