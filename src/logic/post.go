package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

// 帖子相关

func CreatePost(p *models.Post) (err error) {
	// 1.生成post id
	p.ID = snowflake.GenID()

	// 2.保存到数据库
	return mysql.CreatePost(p)
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
