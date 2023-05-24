package main

import (
	"ginchat/controller"
	"ginchat/utils"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	utils.ConfigInit() //这里要用viper的话，路径要修改程当前目录下config的相对路径

	dns := viper.GetString("mysql.dns")
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic("failed to connect mysql")
	}
	// db.AutoMigrate(&controller.User{})  //建立数据表
	//db.AutoMigrate(&controller.UserBasic{}) //
	// db.AutoMigrate(&controller.UserBasic{})

	user := controller.UserBasic{}
	user.Name = "ygwsb12444"
	db.Create(&user) //注意这里传的是指针
}
