package controller

import (
	"ginchat/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Avatar        string //头像
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time `gorm:"column:login_time" json:"login_time"`
	HeartbeatTime time.Time `gorm:"column:heartbeat_time" json:"heartbeat_time"`
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	return data
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).Find(&user)//
	if user.Name == "" {
		return UserBasic{}
	}
	return user
}


func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}
