package service

import "github.com/gin-gonic/gin"

// @Summary ping example
// @Description do ping
// @Tags 首页
// @Accept json
// @Success 200 {string} json{"message"}
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world!(Engineering) ",
	})
}

func ToRegister(c *gin.Context) {
	
}