package router

import (
	"ginchat/docs"
	"ginchat/service"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// // @Summary ping example
// // @Description do ping
// // @Produce json
// // @Success 200 {string} json{}
// // @Router /helloworld [get]
// func Helloworld(g *gin.Context) {
// 	g.JSON(200, "helloworld!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
// }

func Router() {
	r := gin.Default()

	//静态资源引入
	r.Static("/asset", "./asset")
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	r.LoadHTMLGlob("templates/**/*")
	r.LoadHTMLFiles("test_gin_.html")
	//swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//首页
	r.GET("/index", service.GetIndex) //首页
	//r.GET("/", service.GetIndex)      //首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"message": "wowwwwwwwwww!",
		})
	})
	r.GET("/toRegister", service.ToRegister)

	r.Run(":8080")
}
