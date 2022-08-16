package models

// Community 定义社区的结构体
type Community struct {
	ID   int64  `json:"id" db:"community_id"` // 分类比较多的话,可以使用雪花算法来生成
	Name string `json:"name" db:"community_name"`
}
