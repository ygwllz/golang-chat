package controller

import (
	"ginchat/utils"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint //对应的谁 /群 ID
	Type     int  //对应的类型  1好友  2群  3xx
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchUserByGroupId(id uint) []uint {
	contacts := []Contact{}
	res := []uint{}
	utils.DB.Where("target_id = ? and type = 2", id).Find(&contacts)
	for i := range contacts {
		res = append(res, contacts[i].OwnerId)
	}
	return res
}

