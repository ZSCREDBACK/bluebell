package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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

	// 注册业务路由
	r.POST("/signup", controller.RegisterHandler)
	r.POST("/login", controller.LoginHandler)

	// 注册一个测试路由
	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		// JWTAuthMiddleware 用于验证请求头中是否具有有效的JWT,如果不正确则直接在中间件中进行了请求的中断
		v, err := controller.GetCurrentUser(c) // 获取从中间件中传递过来的用户名
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":  "pong",
			"username": v,
		})
	})

	return r
}
