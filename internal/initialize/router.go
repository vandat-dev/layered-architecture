package initialize

import (
	"app/global"
	"app/internal/middlewares"
	"app/internal/routers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine
	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.DisableConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	// middleware - CORS cho tất cả origin (*)
	r.Use(middlewares.CORSMiddleware())

	userRouter := routers.RouterGroupApp.User
	MainGroup := r.Group("/api")
	{
		userRouter.InitUserRouter(MainGroup)

		deliveryRouter := routers.RouterGroupApp.DeliveryFrame
		deliveryRouter.InitScanRouter(MainGroup)
	}

	// Swagger endpoint - với CORS đã được áp dụng
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
