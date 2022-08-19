package models

import "time"

// Community 定义社区的结构体
type Community struct {
	ID   int64  `json:"id" db:"community_id"` // 分类比较多的话,可以使用雪花算法来生成
	Name string `json:"name" db:"community_name"`
}

// CommunityDetail 定义社区详情的结构体
type CommunityDetail struct {
	ID           int64     `json:"id" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"created_time" db:"created_time"`
	// time.Time 类型需要在连接数据库时添加parseTime=True选项
}
