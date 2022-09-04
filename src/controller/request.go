package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

const ContextUserName = "username" // 放在这里避免循环引用
const ContextUserId = "id"

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

// GetCurrentUserId 从请求中获取用户ID
func GetCurrentUserId(c *gin.Context) (id int64, err error) {
	v, ok := c.Get(ContextUserId)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	id, ok = v.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

// GetReqPageSize 获取请求中传递的page和size
func GetReqPageSize(c *gin.Context) (page, size int64) {
	// 1.获取分页参数
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	// 2.解析参数
	var err error
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1 // 没有传递就使用默认值,相当于第几页
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10 // 一页多少条记录
	}

	return page, size
}
