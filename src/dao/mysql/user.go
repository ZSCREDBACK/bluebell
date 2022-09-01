package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"go.uber.org/zap"
)

// 把数据库操作封装成函数,在业务逻辑中调用

// 定义一个加盐字符串
var salt = "bluebell"

// CheckUserExist 判断用户记录是否存在
func CheckUserExist(username string) (err error) {
	// 查询数据库中符合条件的记录数
	sqlStr := "SELECT COUNT(user_id) FROM user WHERE username = ?"
	var count int // 用于存放查询结果

	// 执行查询
	if err = db.Get(&count, sqlStr, username); err != nil {
		zap.L().Error("Check user exist failed", zap.Error(err))
		return
	}

	// 如果count大于0,则等式为true,则用户存在
	if count > 0 {
		return ErrUserExist
	}
	return
}

// InsertUser 插入一条用户记录
func InsertUser(user *models.User) error {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)

	// 定义SQL语句
	sqlStr := "INSERT INTO user(user_id, username, password) VALUES(?, ?, ?)"

	// 执行插入操作
	_, err := db.Exec(sqlStr, user.ID, user.Username, user.Password)
	if err != nil {
		zap.L().Error("Insert new user failed", zap.Error(err))
		return err
	}
	return nil
}

// Login 用户登录验证
func Login(user *models.User) (err error) {
	// 从结构体中获取原密码
	oPassword := user.Password

	// 定义SQL语句
	sqlStr := "SELECT user_id, username, password FROM user WHERE username = ?"

	// 查询用户
	// 注意此处不能新建一个结构体去接收值,必须使用传参进行接收,否则会导致后续函数接收不到userID数据(一直为0)
	err = db.Get(user, sqlStr, user.Username)
	// fmt.Println("判断用户是否存在")
	if err == sql.ErrNoRows { // 判断是否查询到记录(用户是否存在)
		//return ErrUserNotExist
		return ErrInvalidPassword

		// 注意: 一般查询用户不存在,不要返回用户不存在这种信息,返回登录错误即可
		// 避免用户恶意查询数据库
	}
	if err != nil { // 如果查询出错,则返回错误
		zap.L().Error("Query user failed", zap.Error(err))
		return
	}

	// 调试
	// zap.L().Debug("返回数据库查询出来的用户ID:", zap.String("the user_id is", strconv.FormatInt(user.ID, 10)))

	// 判断用户的登录密码是否正确
	if user.Password != encryptPassword(oPassword) { // 将数据库中查询出来的密码与用户输入的密码进行比较
		return ErrInvalidPassword
	}

	return
}

// 密码加盐加密
func encryptPassword(oPassword string) string {
	nPassword := md5.New()                                      // 创建一个md5对象
	nPassword.Write([]byte(salt))                               // 加盐
	return hex.EncodeToString(nPassword.Sum([]byte(oPassword))) // 将密码加密后转换成16进制的字符串
}
