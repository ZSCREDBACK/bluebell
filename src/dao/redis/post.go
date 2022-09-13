package redis

import (
	"bluebell/models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

// CreatePost 发帖后需要在redis中进行的操作
func CreatePost(postID, communityID int64) (err error) {
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

	// 3.把这个帖子的id加入对应的社区的set中
	cKey := getKey(KeyCommunityPrefix + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)

	_, err = pipeline.Exec() // 执行事务
	return
}

// 封装一下公共模块,根据key进行ids查询
func getIDsFormKey(key string, page, size int64) ([]string, error) {
	// 确定redis中查询索引的起始点
	start := (page - 1) * size // 倒数第二页结束为止,所有元素的数量之和 // 等同于最后一页第一个元素的序号值(因为索引都是从0开始的)
	end := start + (size - 1)  // 最后一页第一个元素的序号 + 页数大小 - 1 = 最后一个元素的序号值

	// 查询结果,按分数的从高到低返回一个关于post_id的字符串切片
	return rdb.ZRevRange(key, start, end).Result()
}

// GetPostIdsInOrder 按给定的排序要求,查询符合条件的帖子ID
func GetPostIdsInOrder(p *models.ParamPostList) ([]string, error) {
	// 1.根据用户请求中携带的query string参数order去判断帖子的排序规则,再去对应的key中进行记录查找
	var key string
	if p.Order == models.OrderScore {
		key = getKey(KeyPostScore)
	} else {
		key = getKey(KeyPostTime)
	}

	// 2.获取ids
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据帖子ID列表,去获取每篇帖子的投赞成票情况
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getKey(KeyVotedPrefix + id)
	//	// 查找Key中分数为1的元素数量 -> 统计每个帖子的赞成票票数(反对票就是-1的个数)
	//	v := rdb.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}

	// 上面的方法会导致大量的RTT
	// 可以通过pipeline的方式,对代码进行优化(一次性发送多条命令,让redis直接返回计算后的结果)
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getKey(KeyVotedPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders { // 遍历结果
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}

	return
}

// GetCommunityPostIdsInOrder 根据社区查询符合条件的帖子ID
func GetCommunityPostIdsInOrder(p *models.ParamPostList) ([]string, error) {
	var orderKey string
	if p.Order == models.OrderScore {
		orderKey = getKey(KeyPostScore)
	} else {
		orderKey = getKey(KeyPostTime)
	}

	// 使用 zinterstore 把分区的帖子 set 与帖子分数的 zset 生成一个新的 zset
	// 针对新的 zset 按照之前的方式取数据

	// 生成社区key
	cKey := getKey(KeyCommunityPrefix + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存key减少zinterstore操作的执行次数
	key := orderKey + strconv.Itoa(int(p.CommunityID)) // 取交集后key的名称
	if rdb.Exists(key).Val() < 1 {
		// 如果 key 不存在就需要执行如下聚合操作

		// 事务操作
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Weights:   nil,   // 聚合权重,此处可省略
			Aggregate: "MAX", // 聚合规则: 以最大值进行聚合
		}, cKey, orderKey) // 将cKey和orderKey取交集后,聚合到新的key中
		pipeline.Expire(key, 60*time.Second) // 设置key超时时间(避免60s内的重复查询,但是会造成数据更新不及时)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}

	// 存在就根据key直接查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}
