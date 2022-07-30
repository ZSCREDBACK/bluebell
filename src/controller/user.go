package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
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
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}

		// 将validator.ValidationErrors类型的错误则进行翻译
		c.JSON(http.StatusOK, gin.H{
			//"msg": errs.Translate(trans),
			// 使用removeTopStruct函数去除字段名中的结构体名称标识
			"msg": RemoveTopStruct(errs.Translate(trans)),
		})
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
		c.JSON(http.StatusOK, gin.H{
			"message": "注册失败",
		})
		return
	}

	// 4.返回响应
	c.JSON(200, gin.H{
		"message": "sign_up_ok",
	})
}

// LoginHandler 登录用户
func LoginHandler(c *gin.Context) {
	// 1.获取参数并进行参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {

		zap.L().Error("Bind param failed", zap.Error(err))

		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": RemoveTopStruct(errs.Translate(trans)),
		})
		return
	}

	// 2.业务处理
	if err := logic.Login(p); err != nil {
		// 将登录错误的用户记录到服务日志中
		zap.L().Error("Login failed", zap.String("用户", p.Username), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"message": "登录失败,用户名或密码错误!",
		})
		return
	}

	// 3.返回响应
	c.JSON(200, gin.H{
		"message": "登录成功",
	})
}
