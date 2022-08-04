package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// RegisterHandler 注册用户
func RegisterHandler(c *gin.Context) {
	// 1.获取参数并进行参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil { // ShouldBindJSON只能用于检验字段的类型与格式
		// 请求参数有误,直接返回响应
		zap.L().Error("Bind param failed", zap.Error(err))

		// 使用gin内置的validator进行参数校验

		// 判断是否是validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseErr(c, BadRequest) // 响应错误(请求参数不正确)
			return
		}

		// 将validator.ValidationErrors类型的错误则进行翻译
		//c.JSON(http.StatusOK, gin.H{
		//	//"msg": errs.Translate(trans),
		//	// 使用removeTopStruct函数去除字段名中的结构体名称标识
		//	"msg": RemoveTopStruct(errs.Translate(trans)),
		//})
		ResponseErrWithMsg(c, BadRequest, RemoveTopStruct(errs.Translate(trans)))

		return
	}

	// 2.手动对参数进行详细的业务校验(可选)
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	//	zap.L().Error("Request param correct")
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"message": "请求参数不正确",
	//	})
	//	return
	//}

	// 3.业务处理
	if err := logic.Register(p); err != nil {
		zap.L().Error("Register failed", zap.Error(err))
		if errors.Is(err, mysql.ErrUserExist) { // errors.Is用于判断err是否是mysql.ErrUserExist类型
			ResponseErrWithMsg(c, BadRequest, "注册失败,用户名已存在")
			return
		}
		ResponseErr(c, ServerError)
		return
	}

	// 4.返回响应
	ResponseOk(c, nil)
}

// LoginHandler 登录用户
func LoginHandler(c *gin.Context) {
	// 1.获取参数并进行参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {

		zap.L().Error("Bind param failed", zap.Error(err))

		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseErr(c, BadRequest) // 响应错误(请求参数不正确)
			return
		}

		//c.JSON(http.StatusOK, gin.H{
		//	"msg": RemoveTopStruct(errs.Translate(trans)),
		//})
		ResponseErrWithMsg(c, BadRequest, RemoveTopStruct(errs.Translate(trans)))

		return
	}

	// 2.业务处理
	token, err := logic.Login(p)
	if err != nil {
		// 将登录错误的用户记录到服务日志中
		zap.L().Error("Login failed", zap.String("用户", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrUserNotExist) {
			ResponseErr(c, NotFound) // 用户不存在
			return
		}
		ResponseErr(c, AccountOrPasswordError)
		return
	}

	// 3.返回响应
	ResponseOk(c, token)
}
