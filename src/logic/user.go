package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

// 存放业务逻辑代码

// Register 用户注册
func Register(p *models.ParamSignUp) (err error) {
	// 1.判断用户是否已经存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return
	}

	// 2.生成UID
	userID := snowflake.GenID()

	// 3.构造一个User实例
	user := &models.User{
		ID:       userID,
		Username: p.Username,
		Password: p.Password,
	}

	// 4.保存进数据库
	return mysql.InsertUser(user)
}

// Login 用户登录
func Login(p *models.ParamLogin) error {
	// 1.构造一个User实例
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	// 2.设置用户登录状态(可选)

	// 3.查询用户,并进行密码校验
	return mysql.Login(user)
}
