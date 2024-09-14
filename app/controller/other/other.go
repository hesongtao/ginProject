package other

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Wash(c *gin.Context) {
	// 获取get方式参数
	// username := c.Param("username")
	username := c.Query("username")
	// 获取post方式参数
	// username := c.PostForm("username")
	// 获取header参数
	// c.Request.Header.Get("username")

	data := make(map[string]interface{}, 0)

	data["username"] = username

	c.JSON(http.StatusOK, gin.H{
		"message": "fail",
		"status":  http.StatusOK,
		"data":    data,
	})
	c.Abort()
	return
}
