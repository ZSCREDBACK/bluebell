package models

// ParamSignUp 定义注册请求的结构体参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 定义登录请求的结构体参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	// UserID // 从(发请求的用户的)请求中获取用户id
	PostID int64 `json:"post_id,string" binding:"required"` // 对哪个帖子进行投票
	// Direction int8 `json:"direction,string" binding:"required,oneof=1 0 -1"` // 二选一
	Direction *int8 `json:"direction" binding:"required,oneof=1 0 -1"` // 注意这里是指针
	// 赞成(1)or反对(-1)or取消投票(0)
	// 使用 validator 内置的 oneof 进行校验,使其值必须在预定义的范围内{1,0,-1}
	// 注意如果Direction设置成数值型(而非像PostID一样的字符串)
	// 需要注意反序列化时的零值(即传入参数等于字段类型默认值的情况(0,false,""))问题
	// 零值不会进行传递到结构体中,传递零值反序列化时会对该字段进行忽略,想要解决这个问题,需要使用指针
	// 如果是gorm的话,也可以使用内置的sql.NullInt64进行解决
	// 也可以通过去除required tag来避免(validator库存在的问题),但是个人不推荐这样做
}
