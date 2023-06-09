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

//群里的用户id列表
func SearchUserByGroupId(id uint) []uint {
	contacts := []Contact{}
	res := []uint{}
	utils.DB.Where("target_id = ? and type = 2", id).Find(&contacts)
	for i := range contacts {
		res = append(res, contacts[i].OwnerId)
	}
	return res
}

func AddFriend(userId uint, targetName string) (int, string) {
	if targetName != "" {
		contact := Contact{}
		user := FindUserByName(targetName)
		if user.Name == "" {
			return -1, "用户名无效"
		}
		if user.ID == userId {
			return -1, "不能添加自己"
		}

		utils.DB.Where("owner_id = ? and target_id = ? and type = 1", userId, user.ID).Find(&contact)
		if contact.ID != 0 {
			return -1, "该好友已存在"
		}
		tx := utils.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		contact.OwnerId = userId
		contact.TargetId = user.ID
		contact.Type = 1
		if err := utils.DB.Create(&contact).Error; err != nil {
			tx.Rollback()
			return -1, "添加好友失败"
		}
		contact.OwnerId = user.ID
		contact.TargetId = userId
		contact.Type = 1
		if err := utils.DB.Create(&contact).Error; err != nil {
			tx.Rollback()
			return -1, "添加好友失败"
		}
		tx.Commit()
		return 0, "添加好友成功"

	}
	return -1, "用户名为空"
}

//好友列表
func SearchFriends(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type = 1", userId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}