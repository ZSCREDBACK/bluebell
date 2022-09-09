package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"time"
)

// redis 基础: https://redis.cn/

// 推荐阅读
// 基于用户投票的相关算法: https://www.ruanyifeng.com/blog/algorithm/

// 这里使用简单的投票分数
// 投一票加432分 -> 86400(一天)/200 -> 200票可以给帖子续一天热度(此方法来源于 <<redis实战>>)

/*
投票情况分析
Direction=1:
	1.之前没有投过票 -> 更新分数和投票记录 -> +432
	2.之前投过反对票 -> 更新分数和投票记录 -> +864
Direction=0
	1.投过赞成票,想取消 -> 更新分数和投票记录 -> -432
	2.投过反对票,想取消 -> 更新分数和投票记录 -> +432
Direction=-1
	1.之前没有投过票 -> 更新分数和投票记录 -> -432
	2.之前投过赞成票 -> 更新分数和投票记录 -> -864

投票限制
	1.帖子仅允许自发表之日起一周内进行投票(避免大量唤醒冷数据,降低后端数据库的查询压力)
	2.到期之后,将redis中统计的帖子票数持久化到mysql中
	2.然后删除redis中保存的为该帖子投票时创建的KeyVotedPrefix:post_id
*/

// 定义一周的秒数
const (
	oneWeekOnSec = 7 * 24 * 3600
	scorePerVote = 432 // 每一票的分值
)

var ErrVoteTimeExpire = errors.New("帖子投票时间已过")

func VoteForPost(userID, postID string, value float64) error {
	// 1.判断投票限制(取帖子发布时间)
	postTime := rdb.ZScore(getKey(KeyPostTime), postID).Val() // GO 语言操作 redis ZScore 返回一个float64类型
	if float64(time.Now().Unix())-postTime > oneWeekOnSec {
		return ErrVoteTimeExpire
	}

	// 2.更新帖子分数

	// 2.1.查询当前用户给该帖子的投票记录
	ov := rdb.ZScore(getKey(KeyVotedPrefix+postID), userID).Val() // 之前的分数
	var op float64
	if value > ov { // 如果现在投票的分数大于之前的分数,则认为是赞成票,方向为正
		op = 1
	} else { // 反之是反对票
		op = -1
	}

	// 2.2.计算新分数与旧分数之间的差值的绝对值
	diff := math.Abs(ov - value)

	// 做好准备工作后,开始执行事务
	pipeline := rdb.TxPipeline()

	// 2.3.更新帖子的分数
	pipeline.ZIncrBy(getKey(KeyPostScore), op*diff*scorePerVote, postID)

	// 3.记录用户为该帖子投票的数据
	if value == 0 { // 如果direction为0,就去除该用户对这个帖子的投票分数
		pipeline.ZRem(getKey(KeyVotedPrefix+postID), userID) // 这里应是userID,别写错了?
	} else {
		pipeline.ZAdd(getKey(KeyVotedPrefix+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}

	// 执行事务
	_, err := pipeline.Exec()

	return err
}
