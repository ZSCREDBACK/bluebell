package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

// 帖子相关

func CreatePost(p *models.Post) (err error) {
	// 1.生成post id
	p.ID = snowflake.GenID()

	// 2.保存到数据库
	if err = mysql.CreatePost(p); err != nil {
		return err
	}

	// 3.同步发帖时间到redis,使其成为投票截止时间的依据
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
}

func GetPostById(id int64) (data *models.ApiPostDetail, err error) {
	// 查询帖子详情
	post, err := mysql.GetPostById(id)
	if err != nil {
		zap.L().Error("Get post by id failed",
			zap.Int64("post_id", id),
			zap.Error(err))
		return
	}

	// 根据作者id查询出对应作者
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("Get post's author failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}

	// 根据社区id查询出社区详细信息
	community, err := mysql.GetCommunityDetail(post.CommunityID) // 社区ID必须是已经存在的,否则会报错: 无效的ID
	if err != nil {
		zap.L().Error("Get community detail failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}

	// 组合出我们想要的数据
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}

	return
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	// 分页获取帖子列表
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}

	// 初始化切片
	data = make([]*models.ApiPostDetail, 0, len(postList))

	// 遍历帖子列表
	for _, post := range postList {
		// 根据作者id查询出对应作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("Get post's author failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		// 根据社区id查询出社区详细信息
		community, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("Get community detail failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		// 组合数据
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		// 追加到data切片中
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 1.从redis中查询id列表
	ids, err := redis.GetPostIdsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("Get post list 2 failed, no post id found.", zap.Error(err))
		return
	}

	// 2.根据id去mysql数据库中查询帖子的详细信息,并按照我们给定的顺序返回
	postList, err := mysql.GetPostListByIds(ids)
	if err != nil {
		zap.L().Error("Get post list by ids failed", zap.Error(err))
		return
	}

	// 查询出每个帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	// 3.返回数据(复制上面的即可)
	for idx, post := range postList {
		// 根据作者id查询出对应作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("Get post's author failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		// 根据社区id查询出社区详细信息
		community, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("Get community detail failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		// 组合数据
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx], // 利用idx的数量对该帖子赞成票进行计数
			Post:            post,
			CommunityDetail: community,
		}

		// 追加到data切片中
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetCommunityPostIdsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("Get community post ids in order failed, no post id found.", zap.Error(err))
		return
	}

	postList, err := mysql.GetPostListByIds(ids)
	if err != nil {
		zap.L().Error("Get community post list by ids failed", zap.Error(err))
		return
	}

	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	for idx, post := range postList {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("Get post's author failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		community, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("Get community detail failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postDetail)
	}
	return
}

// GetPostListNew 将两个帖子列表的查询逻辑合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 判断是否传入了社区ID
	if p.CommunityID == 0 {
		return GetPostList2(p) // 查询所有帖子
	} else {
		return GetCommunityPostList(p) // 根据社区ID查询所有帖子
	}
}
