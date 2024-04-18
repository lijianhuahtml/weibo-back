package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// Ok 接口执行正常 需要返回数据 data
func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

// OkWithMsg 接口执行正常 需要返回数据 data
func OkWithMsg(c *gin.Context, message string,
	data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

func Fail(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    -1,
		Message: message,
		Data:    nil,
	})
	c.Abort()
}

func Error(c *gin.Context, message string, code int) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
	c.Abort()
}
