package routes

import (
	"bluebell/logger"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(Mode string) *gin.Engine {
	switch Mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		zap.L().Warn("gin mode unknown: "+Mode+" (available mode: debug release test)",
			zap.String("will set mode in", "debug"))
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册一个测试路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	return r
}
