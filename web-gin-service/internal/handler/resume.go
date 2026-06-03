package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"

	"web-gin-service/internal/grpc"
	"web-gin-service/internal/middleware"
	"web-gin-service/pkg"
	"web-gin-service/proto/gen/resume"
)

type ResumeHandler struct{}

func NewResumeHandler() *ResumeHandler {
	return &ResumeHandler{}
}

func (h *ResumeHandler) GetUploadSignURL(c *gin.Context) {
	candidateID := middleware.GetUserID(c)
	if candidateID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	fileName := c.Query("file_name")
	contentType := c.Query("content_type")

	if fileName == "" {
		pkg.ParamError(c, "文件名不能为空")
		return
	}

	if contentType == "" {
		pkg.ParamError(c, "Content-Type不能为空")
		return
	}

	grpcReq := &resume.GetUploadSignURLReq{
		CandidateId: candidateID,
		FileName:    fileName,
		ContentType: contentType,
	}

	resp, err := grpc.GetResumeClient().GetUploadSignURL(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取上传签名URL失败："+err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"upload_url": resp.UploadUrl,
		"oss_key":    resp.OssKey,
		"expire_sec": resp.ExpireSec,
	})
}

func (h *ResumeHandler) ConfirmUpload(c *gin.Context) {
	candidateID := middleware.GetUserID(c)
	if candidateID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	var req struct {
		OssKey   string `json:"oss_key" binding:"required"`
		FileName string `json:"file_name" binding:"required"`
		FileType string `json:"file_type" binding:"required"`
		FileSize int64  `json:"file_size" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ParamError(c, "参数错误："+err.Error())
		return
	}

	if req.FileType != "pdf" && req.FileType != "doc" && req.FileType != "docx" {
		pkg.ParamError(c, "不支持的文件格式，仅支持 PDF、DOC、DOCX")
		return
	}

	grpcReq := &resume.ConfirmUploadReq{
		CandidateId: candidateID,
		OssKey:      req.OssKey,
		FileName:    req.FileName,
		FileType:    req.FileType,
		FileSize:    req.FileSize,
	}

	resp, err := grpc.GetResumeClient().ConfirmUpload(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "确认上传失败："+err.Error())
		return
	}

	if resp.ResumeId == 0 {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, gin.H{
		"resume_id": resp.ResumeId,
	})
}

func (h *ResumeHandler) GetDownloadSignURL(c *gin.Context) {
	requesterID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	if requesterID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	resumeIDStr := c.Param("id")
	resumeID, err := strconv.ParseInt(resumeIDStr, 10, 64)
	if err != nil {
		pkg.ParamError(c, "无效的简历ID")
		return
	}

	grpcReq := &resume.GetDownloadSignURLReq{
		ResumeId:    resumeID,
		RequesterId: requesterID,
		Role:        role,
	}

	resp, err := grpc.GetResumeClient().GetDownloadSignURL(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取下载签名URL失败："+err.Error())
		return
	}

	if resp.DownloadUrl == "" {
		pkg.Forbidden(c, "无权下载该简历")
		return
	}

	pkg.Success(c, gin.H{
		"download_url": resp.DownloadUrl,
		"file_name":    resp.FileName,
		"expire_sec":   resp.ExpireSec,
	})
}

func (h *ResumeHandler) ListMyResumes(c *gin.Context) {
	candidateID := middleware.GetUserID(c)
	if candidateID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	grpcReq := &resume.ListMyResumesReq{
		CandidateId: candidateID,
	}

	resp, err := grpc.GetResumeClient().ListMyResumes(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取简历列表失败："+err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"resumes": resp.Resumes,
	})
}

func (h *ResumeHandler) DeleteResume(c *gin.Context) {
	candidateID := middleware.GetUserID(c)
	if candidateID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	resumeIDStr := c.Param("id")
	resumeID, err := strconv.ParseInt(resumeIDStr, 10, 64)
	if err != nil {
		pkg.ParamError(c, "无效的简历ID")
		return
	}

	grpcReq := &resume.DeleteResumeReq{
		ResumeId:    resumeID,
		CandidateId: candidateID,
	}

	resp, err := grpc.GetResumeClient().DeleteResume(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "删除失败："+err.Error())
		return
	}

	if !resp.Success {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, nil)
}
