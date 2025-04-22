package handler

import (
	minioHandler "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/handler/minio_handler"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/minio"
	"github.com/gin-gonic/gin"
)

type Services struct {
	minioService minio.Client
}

type Handlers struct {
	minioHandler minioHandler.Handler
}

func NewHandler(minioService minio.Client) (*Services, *Handlers) {
	
	return &Services{minioService: minioService,},
	&Handlers{minioHandler: *minioHandler.NewMinioHandler(minioService)}
}

func (h *Handlers) RegisterRoutes(router *gin.Engine) {

	minioRoutes := router.Group("/files")
	{
		minioRoutes.POST("/one", h.minioHandler.CreateOne)
		minioRoutes.POST("/many", h.minioHandler.CreateMany)

		minioRoutes.GET("/one", h.minioHandler.GetOne)
		minioRoutes.GET("/many", h.minioHandler.GetMany)
		minioRoutes.GET("/all", h.minioHandler.GetAll)

		minioRoutes.DELETE("/one", h.minioHandler.DeleteOne)
		minioRoutes.DELETE("/many", h.minioHandler.DeleteMany)
	}

}
