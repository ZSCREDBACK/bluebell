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
	NeedLogin
	InvalidToken
	UnknownError
	ParamError
)

var responseCodeMsg = map[ResponseCode]string{
	Success:                "成功",
	BadRequest:             "错误的请求",
	Unauthorized:           "未授权",
	Forbidden:              "禁止访问",
	AlreadyExist:           "已存在",
	NotFound:               "未找到",
	ServerError:            "服务端错误",
	AccountOrPasswordError: "账号或密码错误",
	NeedLogin:              "需要登录",
	InvalidToken:           "无效的token",
	UnknownError:           "未知错误",
	ParamError:             "参数错误",
}

func (code ResponseCode) CodeToMsg() interface{} {
	msg, ok := responseCodeMsg[code]
	if !ok {
		return "未找到Code对应的提示信息"
	}
	return msg
}
