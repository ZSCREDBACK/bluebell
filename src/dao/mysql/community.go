package mysql

import (
	"bluebell/models"
	"database/sql"
	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("no community found")
			return nil, nil
		}
	}
	return
}

func GetCommunityDetail(communityId int64) (community *models.CommunityDetail, err error) {
	sqlStr := "select community_id, community_name, introduction, created_time from community where community_id = ?"
	community = new(models.CommunityDetail)
	if err = db.Get(community, sqlStr, communityId); err != nil {
		if err == sql.ErrNoRows {
			// 未找到ID对应的社区
			err = ErrInvalidID
		}
	}
	return community, err
}
