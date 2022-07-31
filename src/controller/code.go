package controller

// 定义返回状态码的含义

type ResponseCode int64

const (
	Success ResponseCode = 1000 + iota
	BadRequest
	Unauthorized
	Forbidden
	AlreadyExist
	NotFound
	ServerError
	AccountOrPasswordError
	UnknownError
)

var responseCodeMsg = map[ResponseCode]string{
	Success:                "success",
	BadRequest:             "请求参数不正确",
	Unauthorized:           "未授权",
	Forbidden:              "拒绝访问",
	AlreadyExist:           "已存在",
	NotFound:               "未找到",
	ServerError:            "服务器错误",
	AccountOrPasswordError: "账号或密码错误",
	UnknownError:           "未知错误",
}

func (code ResponseCode) CodeToMsg() interface{} {
	msg, ok := responseCodeMsg[code]
	if !ok {
		return "未找到Code对应的提示信息"
	}
	return msg
}
