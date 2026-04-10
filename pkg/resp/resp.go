package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// R 统一响应结构
type R struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, R{Code: 0, Msg: "success", Data: data})
}

// Fail 失败响应
// code < 1000 时作为 HTTP 状态码（如 400、401、404、500）
// code >= 1000 时为业务错误码，HTTP 状态码固定 200
func Fail(c *gin.Context, code int, msg string) {
	httpCode := http.StatusOK
	if code > 0 && code < 1000 {
		httpCode = code
	}
	c.JSON(httpCode, R{Code: code, Msg: msg})
}

// Page 分页数据结构（配合 OK 使用）
type Page struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}
