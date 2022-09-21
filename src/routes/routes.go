package routes

import (
	"bluebell/controller"
	_ "bluebell/docs" // 导入生成的docs
	"bluebell/logger"
	"bluebell/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-contrib/pprof"
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
	// 注册全局限流中间件: 两秒填充一个令牌,令牌桶中最多有个1个令牌
	// r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))

	v1 := r.Group("/api/v1")

	// 注册路由
	v1.POST("/signup", controller.RegisterHandler) // 注册
	v1.POST("/login", controller.LoginHandler)     // 登录

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // 接口文档相关
	// r.GET("/swagger/*any", gs.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE")) // 该环境变量如果有值,则禁用swagger

	// 注册认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())

	// 注册一个测试路由(测试用)
	//v1.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	//	// JWTAuthMiddleware 用于验证请求头中是否具有有效的JWT,如果不正确则直接在中间件中进行了请求的中断
	//	v, err := controller.GetCurrentUser(c) // 获取从中间件中传递过来的用户名
	//	if err != nil {
	//		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	//		return
	//	}
	//	c.JSON(http.StatusOK, gin.H{
	//		"message":  "pong",
	//		"username": v,
	//	})
	//})

	// 注册业务路由
	{
		v1.GET("/community", controller.CommunityHandler)           // 获取社区列表
		v1.GET("/community/:id", controller.CommunityDetailHandler) // 获取社区详情
		v1.POST("/post", controller.CreatePostHandler)              // 发帖
		v1.GET("/post/:id", controller.GetPostDetailHandler)        // 获取帖子详情
		v1.GET("/posts", controller.GetPostListHandler)             // 获取帖子列表
		v1.POST("/vote", controller.PostVoteHandler)                // 用户投票
		v1.GET("/posts2", controller.GetPostListHandler2)           // 获取帖子列表(优化版)
	}

	// 注册pprof相关路由,用于性能调优
	// 调用性能分析工具进行CPU分析: go tool pprof http://127.0.0.1:12345/debug/pprof/profile
	// 使用内存分析: go tool pprof -inuse_space http://127.0.0.1:12345/debug/pprof/heap
	// 配合压测工具使用: go-wrk.exe -t=8  -n=10000 http://127.0.0.1:12345/api/v1/posts
	// 生成火焰图: go-torch -u http://127.0.0.1:12345 -t 30
	// 更多详见: Go性能调优.md
	pprof.Register(r)

	// 定义404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 not found"})
	})

	return r
}
