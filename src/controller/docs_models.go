package controller

import "bluebell/models"

// 专门用于存储swagger文档中使用的响应数据models

// _ResponsePostList 帖子列表接口响应数据
type _ResponsePostList struct {
	Code    ResponseCode            `json:"code"` // 业务响应状态码
	Message string                  `json:"msg"`  // 提示信息
	Data    []*models.ApiPostDetail `json:"data"` // 数据
}

// _ResponsePostDetail 帖子详情接口响应数据
type _ResponsePostDetail struct {
	Code    ResponseCode          `json:"code"` // 业务响应状态码
	Message string                `json:"msg"`  // 提示信息
	Data    *models.ApiPostDetail `json:"data"` // 数据
}

// 示例结构体的返回值未嵌套的问题待解决(直接修改最终生成的docs.go?!)
