package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查询数据库中的社区列表并返回
	return mysql.GetCommunityList()
}
