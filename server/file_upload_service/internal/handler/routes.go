package handler

import (
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/services/minio"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/services/quota"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	minioService minio.Client
	quotaService *quota.QuotaService
}

func NewHandler(minioService minio.Client, quotaService *quota.QuotaService) *Handler {
	return &Handler{
		minioService: minioService,
		quotaService: quotaService,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	routesUserApi := router.Group("/user")
	{
		routesUserApi.POST("/:id/plan/init", h.InitUserPlan)
		routesUserApi.GET("/:id/usage", h.GetUserUsage)
	}

	routesFileApi := router.Group("/files")
	{
		routesFileApi.POST("/one", h.CreateOne)
		routesFileApi.POST("/many", h.CreateMany)

		routesFileApi.GET("/one", h.GetOne)
		routesFileApi.GET("/many", h.GetMany)
		routesFileApi.GET("/all", h.GetAll)

		routesFileApi.DELETE("/one", h.DeleteOne)
		routesFileApi.DELETE("/many", h.DeleteMany)
	}

}
