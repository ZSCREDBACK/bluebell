package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

// VoteForPost 为帖子投票
func VoteForPost(userID int64, post *models.ParamVoteData) error {
	// Debug一下投票过程
	zap.L().Debug("Vote for post",
		zap.Int64("user's id", userID),
		zap.Int64("post's id", post.PostID),
		zap.Int8("vote's direction", *post.Direction),
	)

	// 查询 redis 时需要使用 float 类型,这里先进行转换
	// 如果用的频率比较高也可以放在公共模块中
	value, err := strconv.ParseFloat(fmt.Sprintf("%d", post.Direction), 10)
	if err != nil {
		zap.L().Error("Parse direction to float failed", zap.Error(err))
		return err
	}

	return redis.VoteForPost(
		fmt.Sprintf("%d", userID),
		fmt.Sprintf("%d", post.PostID),
		value,
	)
}
