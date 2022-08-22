package model

import (
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
)

/**
 * @Author: rayoluo
 * @Author: hustly123@gmail.com
 * @Date: 2021/5/23 0:24
 * @Desc: 文章和标签对应表
 */

type ArticleTag struct {
	gorm.Model
	ArticleID int
	TagID     int
}

// AddArtTags 添加文章标签对应关系
func AddArtTags(id int, tags []int) int {
	for i := 0; i < len(tags); i++ {
		err := db.Create(&ArticleTag{ArticleID: id, TagID: tags[i]}).Error
		if err != nil {
			return errmsg.ERROR
		}
	}
	return errmsg.SUCCSE
}

// DeleteArtTags 删除指定ID的文章的所有标签
func DeleteByArtID(id int) int {
	err := db.Where("article_id = ?", id).Delete(&ArticleTag{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// DeleteByTagID 根据tag的id删除所有的标签文章对应关系
func DeleteByTagID(id int) int {
	err := db.Where("tag_id = ?", id).Delete(&ArticleTag{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// GetArtTags 根据文章ID查询文章的所有标签，返回taglist captcha total
func GetArtTags(id int) ([]Tag, int, int64) {
	var tagList []Tag
	var total int64
	err := db.Model(&Tag{}).Joins("LEFT JOIN article_tag ON tag.id = article_tag.tag_id").Joins("LEFT JOIN article ON article_tag.article_id = article.id").
		Where("article.id = ?", id).Where("article_tag.deleted_at is null").Find(&tagList).Error
	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	db.Model(&ArticleTag{}).Where("article_id = ?", id).Count(&total)
	return tagList, errmsg.SUCCSE, total
}

// GetArtsByTag 根据标签ID查询标签对应的所有文章，返回article captcha total
func GetArtsByTag(id int, pageSize int, pageNum int) ([]ArtWithTags, int, int64) {
	var artList []Article
	var total int64
	err := db.Model(&Article{}).Preload("Category").Joins("LEFT JOIN article_tag ON article.id = article_tag.article_id").
		Joins("LEFT JOIN tag ON article_tag.tag_id = tag.id").Where("tag.id = ?", id).
		Where("article_tag.deleted_at is null").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&artList).Error
	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	db.Model(&ArticleTag{}).Where("tag_id = ?", id).Count(&total)

	awsl := make([]ArtWithTags, 0)
	for i := 0; i < len(artList); i++ {
		tags, _, _ := GetArtTags(int(artList[i].ID))
		awsl = append(awsl, ArtWithTags{
			artList[i],
			tags,
		})
	}
	return awsl, errmsg.SUCCSE, total
}

type ArtWithTags struct {
	Article
	TagList []Tag `json:"tags"` // 踩坑：小写的不会被json解析
}

// SearchArtTags 根据搜索的文章标题返回"文章+标签"列表
func SearchArtTags(title string, pageSize int, pageNum int) ([]ArtWithTags, int, int64) {
	var (
		arts  []Article
		code  int
		total int64
	)
	arts, code, total = SearchArticle(title, pageSize, pageNum)
	if code == errmsg.ERROR {
		return nil, code, 0
	}
	awsl := make([]ArtWithTags, 0)
	for i := 0; i < len(arts); i++ {
		tags, _, _ := GetArtTags(int(arts[i].ID))
		awsl = append(awsl, ArtWithTags{
			arts[i],
			tags,
		})
	}
	return awsl, errmsg.SUCCSE, total
}

// GetArtTagsList 查询“文章+标签”列表
func GetArtTagsList(pageSize int, pageNum int) ([]ArtWithTags, int, int64) {
	var (
		arts  []Article
		code  int
		total int64
	)
	arts, code, total = GetArt(pageSize, pageNum)
	if code == errmsg.ERROR {
		return nil, code, 0
	}
	awsl := make([]ArtWithTags, 0)
	for i := 0; i < len(arts); i++ {
		tags, _, _ := GetArtTags(int(arts[i].ID))
		awsl = append(awsl, ArtWithTags{
			arts[i],
			tags,
		})
	}
	// log.Println(awsl)
	return awsl, errmsg.SUCCSE, total
}

// GetCateArtTags 获取id指定的分类下的“文章+标签”列表
func GetCateArtTags(id int, pageSize int, pageNum int) ([]ArtWithTags, int, int64) {
	var (
		arts  []Article
		code  int
		total int64
	)
	arts, code, total = GetCateArt(id, pageSize, pageNum)
	if code == errmsg.ERROR {
		return nil, code, 0
	}
	awsl := make([]ArtWithTags, 0)
	for i := 0; i < len(arts); i++ {
		tags, _, _ := GetArtTags(int(arts[i].ID))
		awsl = append(awsl, ArtWithTags{
			arts[i],
			tags,
		})
	}
	return awsl, errmsg.SUCCSE, total
}
