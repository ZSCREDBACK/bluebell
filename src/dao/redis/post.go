package redis

import (
	"github.com/go-redis/redis"
	"time"
)

// CreatePost 发帖后需要在redis中进行的操作
func CreatePost(postID int64) (err error) {
	// 设置事务操作(必须同时成功或失败)
	pipeline := rdb.TxPipeline()

	// 1.设置这个帖子开始投票的时间
	pipeline.ZAdd(getKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 2.设置这个帖子的初始分数
	pipeline.ZAdd(getKey(KeyPostScore), redis.Z{ // 如果报错类型不匹配,就先清理一下redis的key
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	_, err = pipeline.Exec() // 执行事务
	return
}
