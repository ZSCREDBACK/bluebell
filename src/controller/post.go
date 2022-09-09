package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// 帖子相关

func CreatePostHandler(c *gin.Context) {
	// 1.获取参数并进行参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Bind Param Error", zap.Error(err))
		ResponseErr(c, ParamError)
		return
	}

	// 2.从 c 取到当前发请求的用户ID
	userId, err := GetCurrentUserId(c)
	if err != nil {
		ResponseErr(c, NeedLogin)
		return
	}
	p.AuthorID = userId

	// 打印author id排查问题
	// zap.L().Info("打印author id", zap.String("The post's author_id is", strconv.FormatInt(p.AuthorID, 10)))

	// 3.创建帖子
	if err = logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseErr(c, ServerError)
		return
	}

	// 4.返回响应
	ResponseOk(c, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	// 1.获取参数(从URL中获取帖子的id)
	paramStr := c.Param("id")
	parseInt, err := strconv.ParseInt(paramStr, 10, 64)
	if err != nil {
		zap.L().Error("Get post detail with invalid param", zap.Error(err))
		ResponseErr(c, ParamError)
		return
	}

	// 2.根据id获取到帖子的数据(查询数据库)
	data, err := logic.GetPostById(parseInt)
	if err != nil {
		// zap.L().Error("logic.GetPostById failed", zap.Error(err)) // 多余了
		ResponseErr(c, ServerError)
		return
	}

	// 3.返回响应
	ResponseOk(c, data)
}

func GetPostListHandler(c *gin.Context) {
	// 1.获取请求中的分页参数
	page, size := GetReqPageSize(c)

	// 2.获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("Get post list failed", zap.Error(err))
		ResponseErr(c, ServerError)
		return
	}

	// 3.返回数据
	ResponseOk(c, data)
}

func GetPostListHandler2(c *gin.Context) {
	// 1.初始化帖子列表结构体,并设置默认值
	p := &models.ParamPostList{
		Page:  1, // 这里推荐从配置文件中获取
		Size:  10,
		Order: models.OrderTime, // 避免magic string,使用定义常量的方式进行代替
	}

	// GET请求参数(query string): /api/v1/posts?page=3&size=2&order=time
	// 2.从请求中获取帖子的排序方式、页数、分页大小
	// 从请求参数中获取query string参数,也可以使用ShouldBind自动进行识别与绑定
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("Get post list failed, fail to bind the query string param", zap.Error(err))
		ResponseErr(c, ServerError)
		return
	}

	// 3.获取数据
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("Get post list failed", zap.Error(err))
		ResponseErr(c, ServerError)
		return
	}

	// 4.返回数据
	ResponseOk(c, data)
}
