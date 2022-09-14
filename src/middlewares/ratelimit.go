package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/juju/ratelimit"
)

// RateLimitMiddleware 限流中间件,需要传入令牌桶令牌的填充速率和令牌桶最大容量两个参数
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求,并返回 rate limit...
		if bucket.TakeAvailable(1) < 1 { // 判断可用令牌是否小于1
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		// 取到令牌就放行
		c.Next()
	}
}

// 我们可以按照不同的限流策略将其注册到不同的位置
// 1.如果要对全站限流就可以注册成全局的中间件
// 2.如果是某一组路由需要限流,那么就只需将该限流中间件注册到对应的路由组即可
