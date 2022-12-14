package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func PostVoteHandler(c *gin.Context) {
	// 1.参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		err1, ok := err.(validator.ValidationErrors) // 对错误进行类型断言,判断是否属于参数错误
		if !ok {
			ResponseErr(c, ParamError)
			return
		}
		errDate := RemoveTopStruct(err1.Translate(trans)) // 翻译并去除掉错误提示中的结构体标识
		ResponseErrWithMsg(c, ParamError, errDate)
		return
	}

	// 2.获取当前请求用户的ID
	userID, err := GetCurrentUserId(c)
	if err != nil {
		ResponseErr(c, NeedLogin) // 没有获取到userid,说明没有登录
		return
	}

	// 3.进行投票
	err = logic.VoteForPost(userID, p)
	if err != nil {
		zap.L().Error("Vote for post failed", zap.Error(err))
		ResponseErr(c, ServerError)
		return
	}

	// 4.返回响应
	ResponseOk(c, nil)
}
