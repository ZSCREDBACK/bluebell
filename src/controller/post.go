package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// 帖子相关

func CreatePostHandler(c *gin.Context) {
	// 1.获取参数并进行参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Bind Param Error", zap.Error(err))
		ResponseErr(c, ParamError)
		return
	}

	// 2.从 c 取到当前发请求的用户ID
	userID, err := GetCurrentUserId(c)
	if err != nil {
		ResponseErr(c, NeedLogin)
		return
	}
	p.AuthorID = userID

	// 打印author id排查问题
	zap.L().Info("打印author id", zap.String("The post's author_id is", strconv.FormatInt(p.AuthorID, 10)))

	// 3.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseErr(c, ServerError)
		return
	}

	// 4.返回响应
	ResponseOk(c, nil)
}
