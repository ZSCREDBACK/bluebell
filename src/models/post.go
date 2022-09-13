package models

import "time"

// 帖子的数据库模型

type Post struct {
	ID          int64     `json:"id,string" db:"post_id"` // 避免数据失真,合理判断范围大小,并不是所有字段都会超过 1^53 -1
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"created_time" db:"created_time"`
}

// 内存对齐: 相同类型的字段放在一起,减少结构体的整体大小,提高程序性能(减少寻址次数)

type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"`
	VoteNum          int64              `json:"vote_num"`
	*Post                               // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区信息(单独放在一个字段中,避免混合)
}
