package models

// 定义用户结构体

type User struct {
	ID       int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Token    string
}
