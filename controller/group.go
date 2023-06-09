package controller

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

func CreateGroup(group Group) (int, string) {
	tx := utils.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if len(group.Name) == 0 {
		return -1, "群名称不能为空"
	}
	if group.OwnerId == 0 {
		return -1, "请先登录"
	}
	if err := utils.DB.Create(&group).Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return -1, "建群失败"
	}
	contant := Contact{}
	contant.OwnerId = group.OwnerId
	contant.ID = group.ID
	contant.Type = 2
	if err := utils.DB.Create(&contant).Error; err != nil {
		tx.Rollback()
		return -1, "关系添加失败"
	}
	tx.Commit()
	return 0, "创建成功"
}

func LoadCommunity(ownerId uint) ([]*Group, string) { //这里的返回值为什么是group指针
	contants := []Contact{}
	groupid := []uint{}
	res := []*Group{}
	utils.DB.Where("owner_id = ? and type = 2", ownerId).Find(&contants)
	for i := range contants {
		id := contants[i].TargetId
		groupid = append(groupid, id)
	}
	utils.DB.Where("ID in ?", groupid).Find(&res)
	for _, v := range res {
		fmt.Println(v)
	}
	return res, "查询成功"
}

func JoinGroup(userId uint, comId string) (int, string) {
	contact := Contact{}
	contact.OwnerId = userId
	contact.Type = 2
	group := Group{}
	utils.DB.Where("id = ?", comId).First(&group)
	if group.Name == "" {
		return -1, "群号错误"
	}
	utils.DB.Where("owner_id = ? and target_id = ?", userId, comId).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "已在群里"
	}
	tx := utils.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	contact.TargetId = group.ID
	if err := utils.DB.Create(&contact).Error; err != nil {
		tx.Rollback()
		return -1, "添加失败"
	}
	tx.Commit()
	return 0, "加群成功"
}
