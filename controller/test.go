package controller

import "github.com/gin-gonic/gin"

func TestRespon(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"msg":  "测试url",
		"data": nil,
	})
}
