package middlewares

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 'Authorization: Bearer your_secret_token'
		// 这里的具体实现方式要依据你的实际业务情况决定,自己灵活修改下就好了
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2003,
			//	"msg":  "请求头中auth为空",
			//})

			// 使用内置的错误处理函数来处理错误
			controller.ResponseErrWithMsg(c, controller.NeedLogin, "请求头中Authorization内容为空")

			c.Abort()
			return
		}

		// 按空格分割请求中Authorization字段的值
		parts := strings.SplitN(authHeader, " ", 2)
		// 判断Authorization字段是否合法
		if !(len(parts) == 2 && parts[0] == "Bearer") { // 以空格分割后,Authorization中只有两个元素且第一个元素为Bearer才认为是合法的
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2004,
			//	"msg":  "请求头中auth格式有误",
			//})

			// 使用内置的错误处理函数来处理错误
			controller.ResponseErrWithMsg(c, controller.InvalidToken, "请求头中Authorization内容的格式有误")

			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString,我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2005,
			//	"msg":  "无效的Token",
			//})

			// 使用内置的错误处理函数来处理错误
			controller.ResponseErr(c, controller.InvalidToken)

			c.Abort()
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		c.Set(controller.ContextUserName, mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get(controller.ContextUserName)来获取当前请求的用户信息,也可以单独用一个函数封装一下
	}
}
