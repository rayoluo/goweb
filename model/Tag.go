package model

import (
	"ginblog/utils/errmsg"
)

/**
 * @Author: rayoluo
 * @Author: hustly123@gmail.com
 * @Date: 2021/5/18 15:21
 * @Desc: 添加分类管理
 */

type Tag struct {
	ID   uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}

// CheckTagUsed 查询标签是否已经使用
func CheckTagUsed(name string) int {
	var tag Tag
	db.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return errmsg.ERROR_TAG_USED
	}
	return errmsg.SUCCSE
}

// CreateTag 新增标签
func CreateTag(data *Tag) int {
	err := db.Create(data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// GetTagInfo 查询单个标签信息
func GetTagInfo(id int) (Tag, int) {
	var tag Tag
	db.Where("id = ?", id).First(&tag)
	return tag, errmsg.SUCCSE
}

// GetTagList 查询标签列表
func GetTagList(pageSize int, pageNum int) ([]Tag, int64) {
	var tags []Tag
	var total int64
	err = db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&tags).Error
	db.Model(&Tag{}).Count(&total)
	if err != nil {
		return nil, 0
	}
	return tags, total
}

// GetAllTags 获取所有的标签
func GetAllTags() ([]Tag, int64) {
	var tags []Tag
	var total int64
	err := db.Find(&tags).Error
	db.Model(&Tag{}).Count(&total)
	if err != nil {
		return nil, 0
	}
	return tags, total
}

// EditTag 编辑标签信息
func EditTag(id int, data *Tag) int {
	var tag Tag
	var maps = make(map[string]interface{})
	maps["name"] = data.Name

	err = db.Model(&tag).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// DeleteTag 删除标签
func DeleteTag(id int) int {
	var tag Tag
	err = db.Where("id = ? ", id).Delete(&tag).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
