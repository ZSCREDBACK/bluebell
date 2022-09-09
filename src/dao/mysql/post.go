package mysql

import (
	"bluebell/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := "insert into `post` (post_id, title, content, author_id, community_id) values (?, ?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostById 根据post_id获取单个帖子详细信息
func GetPostById(id int64) (data *models.Post, err error) {
	data = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, created_time from post where post_id = ?`
	err = db.Get(data, sqlStr, id)
	return
}

// GetPostList 根据分页查询条件,获取匹配的帖子详细信息列表
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	posts = make([]*models.Post, 0, 10) // 初始化一个合适的容量
	// 进行分页,避免展示太多
	// page: 第一个问号表示偏移量(offset),从什么地方开始取(第几页)
	// size: 第二个问好表示数量(limit),限制取多少个(每一页有多少条记录)
	// 以帖子的创建时间降序排序,实现由新到旧展示
	sqlStr := `select post_id, title, content, author_id, community_id, created_time from post order by created_time desc limit ?,?`
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIds 根据给定的id列表查询帖子数据
func GetPostListByIds(ids []string) (posts []*models.Post, err error) {
	// sqlx参考文档: https://www.liwenzhou.com/posts/Go/sqlx/
	sqlStr := `select post_id, title, content, author_id, community_id, created_time
		from post
		where post_id in (?)
		order by FIND_IN_SET(post_id,?)`
	// FIND_IN_SET 是mysql服务端根据给定的顺序,对查询出来的数据进行排序后,并返回给客户端
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ",")) // 该函数生成一个构造后的sql语句
	// strings.Join(ids, ",") // post_id_1,post_id_2,... // 自行构造集合列表
	if err != nil {
		return
	}

	query = db.Rebind(query)                // 将查询结果重新绑定到变量中
	err = db.Select(&posts, query, args...) // 将sql语句和参数一并发给数据库服务端,并将结果反射到posts中
	// 必须要加...,它会自动帮我们进行参数的切分与追加
	return
}
