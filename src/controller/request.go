package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
)

const ContextUserName = "username" // 放在这里避免循环引用

var ErrorUserNotLogin = errors.New("用户未登录")

// GetCurrentUser 从请求中获取用户名信息
func GetCurrentUser(c *gin.Context) (username string, err error) {
	value, ok := c.Get(ContextUserName)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	username, ok = value.(string)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}
