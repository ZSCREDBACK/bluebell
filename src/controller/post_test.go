package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// post
// 执行单元测试: go test [./] [-v]

func TestCreatePostHandler(t *testing.T) {
	// 由于模块拆分较细,为了避免可能出现的循环引用的问题,这里就自己构造一个路由,而非使用 router.Setup

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	// 正常测试
	body := `{
		"community_id": 1,
		"title": "WOW",
		"content": "just a test"
	}`

	// 尝试注入错误
	//body := `{
	//	"title": "WOW",
	//	"content": "just a test"
	//}`

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body))) // 构造请求
	// http.MethodPost: 等同于 "POST"
	// bytes.NewReader(): 将字符串强转的字节切片,转换成一个Reader对象

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 返回状态码
	assert.Equal(t, 200, w.Code)

	// 判断返回结果是否符合预期
	// 由于请求中不带有access token,所以这个用例会卡在第二步
	// 相当于该用例仅仅是用于测试 第一步/参数校验 这个判断的正确性

	// 方式1: 判断是否包含指定字符串
	// assert.Contains(t, w.Body.String(), "需要登录") // 第二步报错提示是 NeedLogin

	// 方式2: 将响应的内容反序列化到Response,然后判断字段是否符合预期
	res := new(Response)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("Json unmarshal w.Body failed,err: %v\n", err)
	}
	assert.Equal(t, res.Code, NeedLogin)
}
