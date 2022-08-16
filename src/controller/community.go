package controller

import (
	"bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区相关

func CommunityHandler(c *gin.Context) {
	// 获取社区列表(community_id, community_name),以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("get community list failed", zap.Error(err))
		ResponseErr(c, ServerError) // 不轻易把服务端报错返回给客户端
		return
	}
	ResponseOk(c, data)
}
