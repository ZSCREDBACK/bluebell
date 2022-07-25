package routes

import (
	"goScaffold/logger"
	"goScaffold/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *settings.AppConfig) *gin.Engine {
	gin.SetMode(cfg.Mode)

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
