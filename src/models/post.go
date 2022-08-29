package models

import "time"

// 帖子的数据库模型

type Post struct {
	ID          int64     `json:"id" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"created_time" db:"created_time"`
}

// 内存对齐: 相同类型的字段放在一起,减少结构体的整体大小,提高程序性能(减少寻址次数)