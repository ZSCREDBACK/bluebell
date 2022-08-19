package controller

import (
	"bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
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

func CommunityDetailHandler(c *gin.Context) {
	// 获取社区详情(community_id, community_name, introduction, created_time)

	// 获取社区id
	communityId := c.Param("id")
	id, err := strconv.ParseInt(communityId, 10, 64)
	if err != nil {
		zap.L().Error("parse community id failed", zap.Error(err))
		ResponseErr(c, ParamError)
		return
	}

	// 通过id获取该社区详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("get community detail failed", zap.Error(err))
		ResponseErr(c, ServerError)
		return
	}
	ResponseOk(c, data)
}
