package service

import (
	"fmt"
	"ginchat/controller"
	"ginchat/utils"
	"html/template"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary ping example
// @Description do ping
// @Tags 首页
// @Accept json
// @Success 200 {string} json{"message"}
// @Router /index [get]
func GetIndex(c *gin.Context) {
	res, err := template.ParseFiles("index.html", "templates/chat/head.html")
	if err != nil {
		panic(err)
	}
	res.Execute(c.Writer, "index")
}

func ToRegister(c *gin.Context) {
	res, err := template.ParseFiles("templates/user/register.html")
	if err != nil {
		panic(err)
	}
	res.Execute(c.Writer, "register")
}

func CreateUser(c *gin.Context) {
	//get参数
	user := controller.UserBasic{}
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	repassword := c.Request.FormValue("Identity")
	fmt.Println(user, repassword)
	//校验数据合法性
	if user.PassWord == "" || user.Name == "" || repassword == "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "用户名或密码不能为空！",
			"data":    user,
		})
		return
	}
	if user.PassWord != repassword {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "两次密码不一致！",
			"data":    user,
		})
		return
	}
	//查找数据库是否存在
	data := controller.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名已注册！",
			"data":    user,
		})
	}

	//创建用户
	salt := fmt.Sprintf("%06d", rand.Int31()) //给密码加盐
	user.Salt = salt
	user.PassWord = utils.MakePassword(user.PassWord, salt)
	user.LoginTime = time.Now()
	user.LoginOutTime = time.Now()
	user.HeartbeatTime = time.Now()
	controller.CreateUser(user) //应在controller层实现
	c.JSON(200, gin.H{
		"code":    0,
		"message": "用户注册成功！",
		"data":    user,
	})
}

func FindUserByNameAndPwd(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	fmt.Println(name, password)
	//查库
	user := controller.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "该用户不存在",
			"data":    "该用户不存在",
		})
		return
	}

	//密码校验
	ok := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !ok {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "密码错误",
			"data":    "密码错误",
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data := controller.FindUserByNameAndPwd(name, pwd)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    data,
	})

}

func ToChat(c *gin.Context) {
	res, err := template.ParseFiles("templates/chat/index.html",
		"templates/chat/head.html",
		"templates/chat/foot.html",
		"templates/chat/tabmenu.html",
		"templates/chat/concat.html",
		"templates/chat/group.html",
		"templates/chat/profile.html",
		"templates/chat/createcom.html",
		"templates/chat/userinfo.html",
		"templates/chat/main.html")
	if err != nil {
		panic(err)
	}
	userid, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		panic(err)
	}
	token := c.Query("token")
	user := controller.UserBasic{}
	user.ID = uint(userid)
	user.Identity = token

	res.Execute(c.Writer, user)
}

func Chat(c *gin.Context) {
	controller.Chat(c.Writer, c.Request)
}
