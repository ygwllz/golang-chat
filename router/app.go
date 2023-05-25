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
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")

	r.LoadHTMLGlob("templates/**/*")
	//swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// r.GET("/helloworld", Helloworld)
	r.GET("/index", service.GetIndex) //回调函数
	r.GET("/", service.GetIndex)

	r.Run(":8080")
}
