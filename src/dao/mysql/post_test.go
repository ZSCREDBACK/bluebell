package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

// 需要先进行初始化数据库连接,否则会报错空指针
// 必须叫init
func init() {
	cfg := settings.MysqlConfig{
		Host:         "192.168.118.138",
		Port:         3306,
		User:         "root",
		Password:     "123456",
		Dbname:       "bluebell",
		MaxOpenConns: 100,
		MaxIdleConns: 50,
	}

	err := Init(&cfg)
	if err != nil {
		panic(err)
	}
}

// TestCreatePost 测试在数据库中创建帖子
func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          1231231231,
		AuthorID:    1231231231,
		CommunityID: 999,
		Title:       "test",
		Content:     "just a test",
	}

	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("Insert the post into database failed, err: %v\n", err)
	}
	t.Log("Insert the post into database failed")
}

// 注意: 该数据会真实插入数据库;所以在完成单元测试后,请执行清理,或者重新初始化数据库
