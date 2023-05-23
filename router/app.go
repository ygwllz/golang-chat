package router

import (
	"ginchat/service"

	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()
	r.GET("/index", service.GetIndex) //回调函数
	r.GET("/", service.GetIndex)

	r.Run(":8080")
}
