package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 封装并统一各种响应
// 1.简化代码
// 2.统一响应格式
// 3.前后端提高协作效率

// 设计目的
/*
{
	"code": 701,		// 程序返回的状态码
	"data": {},			// 程序返回的具体数据
	"msg": "success",	// 程序返回的提示信息
}
*/

type Response struct {
	Code ResponseCode `json:"code"`
	Msg  interface{}  `json:"msg"`
	Data interface{}  `json:"data,omitempty"` // omitempty 没有数据就忽略该字段的返回
}

// ResponseErr 根据状态码返回相应的错误提示信息
func ResponseErr(c *gin.Context, code ResponseCode) {
	c.JSON(http.StatusOK, &Response{
		Code: code,
		Msg:  code.CodeToMsg(),
	})
}

// ResponseErrWithMsg 返回自定义的错误提示信息
func ResponseErrWithMsg(c *gin.Context, code ResponseCode, msg interface{}) {
	c.JSON(http.StatusOK, &Response{
		Code: code,
		Msg:  msg,
	})
}

// ResponseOk 根据状态码返回相应的成功提示信息
func ResponseOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		Code: Success,
		Msg:  Success.CodeToMsg(),
		Data: data,
	})
}
