package model

import (
	"ginblog/utils/errmsg"
)

type Profile struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"type:varchar(20)" json:"name"`
	Desc      string `gorm:"type:varchar(200)" json:"desc"`
	Qqchat    string `gorm:"type:varchar(200)" json:"qq_chat"`
	Wechat    string `gorm:"type:varchar(100)" json:"wechat"`
	Weibo     string `gorm:"type:varchar(200)" json:"weibo"`
	Bili      string `gorm:"type:varchar(200)" json:"bili"`
	Email     string `gorm:"type:varchar(200)" json:"email"`
	Img       string `gorm:"type:varchar(200)" json:"img"`
	Avatar    string `gorm:"type:varchar(200)" json:"avatar"`
	IcpRecord string `gorm:"type:varchar(200)" json:"icp_record"`
	AboutInfo string `gorm:"type:longtext" json:"aboutInfo"`
}

// GetProfile 获取个人信息设置
func GetProfile(id int) (Profile, int) {
	var profile Profile
	err = db.Where("ID = ?", id).First(&profile).Error
	if err != nil {
		return profile, errmsg.ERROR
	}
	return profile, errmsg.SUCCSE
}

// CheckProfile 判断个人信息是否存在
func CheckProfile(id int) bool {
	var profile Profile
	db.Where("id=?", id).First(&profile)
	if profile.ID > 0 {
		return true
	}
	return false
}

// UpdateProfile 更新个人信息设置
func UpdateProfile(id int, data *Profile) int {
	var profile Profile
	// 首先判断ID是否存在 如果不存在则创建
	if !CheckProfile(id) {
		err = db.Create(&data).Error
		if err != nil {
			return errmsg.ERROR
		}
		return errmsg.SUCCSE
	}
	err = db.Model(&profile).Where("ID = ?", id).Updates(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
