package redis

import (
	"bluebell/models"
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

func GetPostIdsInOrder(p *models.ParamPostList) ([]string, error) {
	// 1.根据用户请求中携带的query string参数order去判断帖子的排序规则,再去对应的key中进行记录查找
	var key string
	if p.Order == models.OrderScore {
		key = getKey(KeyPostScore)
	} else {
		key = getKey(KeyPostTime)
	}

	// 2.确定redis中查询索引的起始点
	start := (p.Page - 1) * p.Size // 倒数第二页结束为止,所有元素的数量之和 // 等同于最后一页第一个元素的序号值(因为索引都是从0开始的)
	end := start + (p.Size - 1)    // 最后一页第一个元素的序号 + 页数大小 - 1 = 最后一个元素的序号值

	// 3.查询结果,按分数的从高到低返回一个关于post_id的字符串切片
	return rdb.ZRevRange(key, start, end).Result()
}
