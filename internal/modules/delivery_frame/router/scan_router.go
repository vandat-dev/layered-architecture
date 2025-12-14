package router

import (
	"app/global"
	"app/internal/modules/delivery_frame/controller"
	"app/internal/modules/delivery_frame/repo"
	"app/internal/modules/delivery_frame/service"
	"app/internal/third_party/redis"
	"app/internal/third_party/s3"

	"github.com/gin-gonic/gin"
)

type ScanRouter struct{}

func (sr *ScanRouter) InitScanRouter(Router *gin.RouterGroup) {
	// Dependency Injection (Manual for now, mimicking wire)
	// In a real scenario with wire, this would be: scanController, _ := wire.InitScanRouterHandler()

	s3Provider := s3.NewS3Provider()
	redisProvider := redis.NewRedisProvider()
	scanRepo := repo.NewScanRepository(global.Postgres)
	scanService := service.NewScanService(scanRepo, s3Provider, redisProvider)
	scanController := controller.NewScanController(scanService)

	group := Router.Group("/delivery-frame")
	{
		group.POST("/scan", scanController.CreateScan)
		group.POST("/upload", scanController.UploadImage)
		group.GET("/images", scanController.GetImages)
		group.DELETE("/folder", scanController.DeleteFolder)
	}
}
