package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// 帖子相关

// 使用swagger生成接口文档,注释一般位于controller层

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

// GetPostDetailHandler 帖子详情接口
// @Summary     帖子详情接口
// @Description 可获取帖子详情接口
// @Tags        帖子查询相关接口
// @Accept      application/json
// @Produce     application/json
// @Param       Authorization header string false "Bearer 用户令牌"
// @Param       id            path   int  true  "查询参数"
// @Security    ApiKeyAuth
// @Success     200 {object} _ResponsePostDetail
// @Router      /post/{id} [get]
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

// GetPostListHandler 帖子列表接口
// @Summary     帖子列表接口
// @Description 可按分页参数查询帖子列表接口
// @Tags        帖子查询相关接口
// @Accept      application/json
// @Produce     application/json
// @Param       Authorization header string false "Bearer 用户令牌"
// @Param       page          query  string false "分页页数"
// @Param       size          query  string false "分页大小"
// @Security    ApiKeyAuth
// @Success     200 {object} _ResponsePostList
// @Router      /posts [get]
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

// GetPostListHandler2 升级版帖子列表接口
// @Summary     升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags        帖子查询相关接口
// @Accept      application/json
// @Produce     application/json
// @Param       Authorization header string               false "Bearer 用户令牌"
// @Param       object        query  models.ParamPostList false "查询参数"
// @Security    ApiKeyAuth
// @Success     200 {object} _ResponsePostList
// @Router      /posts2 [get]
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
	//data, err := logic.GetPostList2(p)
	data, err := logic.GetPostListNew(p) // 更新: 函数糅合
	if err != nil {
		zap.L().Error("Get post list failed", zap.Error(err))
		ResponseErr(c, ServerError)
		return
	}

	// 4.返回数据
	ResponseOk(c, data)
}

// GetCommunityPostListHandler 处理根据社区ID查询帖子列表(此处将两个函数进行糅合,所以进行了注释)
//func GetCommunityPostListHandler(c *gin.Context) {
//	p := &models.ParamPostList{
//		CommunityID: 0,
//		Page:        1,
//		Size:        10,
//		Order:       models.OrderTime,
//	}
//
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("Get community post list failed, fail to bind the query string param", zap.Error(err))
//		ResponseErr(c, ServerError)
//		return
//	}
//
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("Get community post list failed", zap.Error(err))
//		ResponseErr(c, ServerError)
//		return
//	}
//
//	ResponseOk(c, data)
//}

// 函数糅合适用于两个函数实现的功能类似,省去一些重复的步骤
