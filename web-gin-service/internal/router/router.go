package router

import (
	"github.com/gin-gonic/gin"

	"web-gin-service/internal/handler"
	"web-gin-service/internal/middleware"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	return &Router{
		engine: gin.Default(),
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	positionHandler := handler.NewPositionHandler()
	applicationHandler := handler.NewApplicationHandler()
	resumeHandler := handler.NewResumeHandler()
	aiChatHandler := handler.NewAIChatHandler()

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	public := r.engine.Group("/api")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.GET("/positions", positionHandler.ListPositions)
		public.GET("/positions/:id", positionHandler.GetPosition)
	}

	hr := r.engine.Group("/api/hr")
	hr.Use(middleware.JWTMiddleware(), middleware.HRMiddleware())
	{
		hr.POST("/positions", positionHandler.CreatePosition)
		hr.PUT("/positions/:id", positionHandler.UpdatePosition)
		hr.POST("/positions/:id/offline", positionHandler.OfflinePosition)
		hr.POST("/positions/:id/online", positionHandler.OnlinePosition)
		hr.GET("/my-positions", positionHandler.ListMyPositions)
		hr.GET("/positions/:id/candidates", applicationHandler.ListCandidates)
		hr.GET("/candidates/:id", applicationHandler.GetCandidateDetail)
		hr.PUT("/applications/:id/status", applicationHandler.UpdateApplicationStatus)
		hr.GET("/resumes/:id/download-url", resumeHandler.GetDownloadSignURL)
		hr.POST("/ai/chat", aiChatHandler.Chat)
		hr.POST("/ai/chat/stream", aiChatHandler.ChatStream)
		hr.GET("/ai/history", aiChatHandler.GetHistory)
		hr.GET("/ai/sessions", aiChatHandler.ListSessions)
		hr.GET("/ai/stats", aiChatHandler.GetStats)
		hr.DELETE("/ai/sessions/:session_id", aiChatHandler.DeleteSession)
		hr.DELETE("/ai/messages/:message_id", aiChatHandler.DeleteMessage)
	}

	candidate := r.engine.Group("/api/candidate")
	candidate.Use(middleware.JWTMiddleware(), middleware.CandidateMiddleware())
	{
		candidate.GET("/profile", userHandler.GetProfile)
		candidate.PUT("/profile", userHandler.UpdateProfile)
		candidate.POST("/positions/:id/apply", applicationHandler.Apply)
		candidate.GET("/applications", applicationHandler.ListMyApplications)
		candidate.DELETE("/applications/:id", applicationHandler.WithdrawApplication)
		candidate.GET("/resumes/upload-url", resumeHandler.GetUploadSignURL)
		candidate.POST("/resumes/confirm", resumeHandler.ConfirmUpload)
		candidate.GET("/resumes", resumeHandler.ListMyResumes)
		candidate.GET("/resumes/:id/download-url", resumeHandler.GetDownloadSignURL)
		candidate.DELETE("/resumes/:id", resumeHandler.DeleteResume)
	}

	return r.engine
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
