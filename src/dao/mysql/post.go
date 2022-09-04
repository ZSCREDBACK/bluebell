package mysql

import (
	"bluebell/models"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := "insert into `post` (post_id, title, content, author_id, community_id) values (?, ?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostById(id int64) (data *models.Post, err error) {
	data = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, created_time from post where post_id = ?`
	err = db.Get(data, sqlStr, id)
	return
}

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	posts = make([]*models.Post, 0, 2)
	// 进行分页,避免展示太多
	// page: 第一个问号表示偏移量(offset),从什么地方开始取(第几页)
	// size: 第二个问好表示数量(limit),限制取多少个(每一页有多少条记录)
	sqlStr := `select post_id, title, content, author_id, community_id, created_time from post limit ?,?`
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}
