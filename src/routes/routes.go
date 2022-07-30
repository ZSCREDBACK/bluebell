package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Setup(GinMode string) *gin.Engine {
	switch GinMode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		zap.L().Warn("gin mode unknown: "+GinMode+" (available mode: debug release test)",
			zap.String("will set mode in", "debug"))
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册一个测试路由
	//r.GET("/", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "ok",
	//	})
	//})

	// 注册业务路由
	r.POST("/signup", controller.RegisterHandler)
	r.POST("/login", controller.LoginHandler)

	return r
}
