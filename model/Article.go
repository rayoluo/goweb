package model

import (
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
	"log"
)

type Article struct {
	Category Category `gorm:"foreignkey:Cid"`
	gorm.Model
	Title   string `gorm:"type:varchar(100);not null" json:"title"`
	Cid     int    `gorm:"type:int;not null" json:"cid"`
	Desc    string `gorm:"type:varchar(200)" json:"desc"`
	Content string `gorm:"type:longtext" json:"content"`
	Img     string `gorm:"type:varchar(100)" json:"img"`
	// 统计评论数量
	CommentCount int `gorm:"type:int;not null;default:0" json:"comment_count"`
	// 统计阅读量
	ReadCount int `gorm:"type:int;not null;default:0" json:"read_count"`
}

// CreateArt 新增文章
func CreateArt(data *Article) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}

// GetCateArt 查询分类下的所有文章
func GetCateArt(id int, pageSize int, pageNum int) ([]Article, int, int64) {
	var cateArtList []Article
	var total int64
	// preload指定主表中的字段名，执行两条sql语句，并将查询category表的结果（限制id）赋值到Article结构体中
	err = db.Preload("Category").Limit(pageSize).Offset((pageNum-1)*pageSize).Where(
		"cid =?", id).Find(&cateArtList).Error
	db.Model(&cateArtList).Where("cid =?", id).Count(&total)
	if err != nil {
		return nil, errmsg.ERROR_CATE_NOT_EXIST, 0
	}
	return cateArtList, errmsg.SUCCSE, total
}

// GetArtInfo 前端查询单个文章
func GetArtInfo(id int) (Article, int) {
	var art Article
	err = db.Where("id = ?", id).Preload("Category").First(&art).Error
	db.Model(&art).Where("id = ?", id).UpdateColumn("read_count", gorm.Expr("read_count + ?", 1))
	if err != nil {
		return art, errmsg.ERROR_ART_NOT_EXIST
	}
	return art, errmsg.SUCCSE
}

// GetArtInfoEdit 后端查询单个文章信息，不用对阅读量进行更新
func GetArtInfoEdit(id int) (Article, int) {
	var art Article
	err = db.Where("id = ?", id).Preload("Category").First(&art).Error
	// db.Model(&art).Where("id = ?", id).UpdateColumn("read_count", gorm.Expr("read_count + ?", 1))
	if err != nil {
		return art, errmsg.ERROR_ART_NOT_EXIST
	}
	return art, errmsg.SUCCSE
}

// GetArt 查询文章列表
func GetArt(pageSize int, pageNum int) ([]Article, int, int64) {
	var articleList []Article
	var err error
	var total int64

	err = db.Debug().Select("article.id, title, img, content, created_at, updated_at, `desc`, comment_count, read_count").Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("Created_At DESC").Joins("Category").Find(&articleList).Error
	// db.Debug().Preload("Category").Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("Created_At DESC").Find(&articleList)
	// log.Println(articleList)
	// 统计文章表中总共有多少篇文章
	// db.Model(&articleList).Count(&total)
	db.Model(&Article{}).Count(&total)
	// log.Println(total)
	if err != nil {
		log.Println(err)
		return nil, errmsg.ERROR, 0
	}
	return articleList, errmsg.SUCCSE, total
}

// SearchArticle 搜索文章标题
func SearchArticle(title string, pageSize int, pageNum int) ([]Article, int, int64) {
	var articleList []Article
	var err error
	var total int64
	err = db.Select("article.id,title, img, content, created_at, updated_at, `desc`, comment_count, read_count").Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Joins("Category").Where("title LIKE ?",
		title+"%",
	).Find(&articleList).Count(&total).Error
	// 单独计数
	//db.Model(&articleList).Count(&total)

	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	return articleList, errmsg.SUCCSE, total
}

// EditArt 编辑文章
func EditArt(id int, data *Article) int {
	var maps = make(map[string]interface{})
	maps["title"] = data.Title
	maps["cid"] = data.Cid
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["img"] = data.Img

	err = db.Model(&Article{}).Where("id = ? ", id).Updates(&maps).Error
	if err != nil {
		// log.Println("fail")
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// func EditArt_new(id int, data *Article) int {
// 	var art Article
// 	var mp map[string]interface{}
// 	mp["title"] = data.Title
// 	mp["cid"] = data.Cid
// 	mp["desc"] = data.Desc
// 	mp["content"] = data.Content
// 	mp["img"] = data.Img
//
// 	err := db.Model(&art).Where("id = ?", id).Updates(mp).Error
// 	if err != nil {
// 		return errmsg.ERROR
// 	}
// 	return errmsg.SUCCSE
// }

// DeleteArt 删除文章
func DeleteArt(id int) int {
	var art Article
	err = db.Where("id = ? ", id).Delete(&art).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// GetAllArt 获取所有文章
func GetAllArt() ([]Article, int, int64) {
	var articleList []Article
	var err error
	var total int64

	err = db.Select("article.id, title, img, created_at, updated_at, `desc`").Order("Created_At DESC").Find(&articleList).Error
	// db.Debug().Preload("Category").Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("Created_At DESC").Find(&articleList)
	// log.Println(articleList)
	// 统计文章表中总共有多少篇文章
	// db.Model(&articleList).Count(&total)
	db.Model(&Article{}).Count(&total)
	// log.Println(total)
	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	return articleList, errmsg.SUCCSE, total
}

// 根据文章id获取推荐文章
func GetRecArts(id int) ([]Article, int, int64) {
	var (
		articleTagList, atls []ArticleTag
		articleList          []Article
		tagIDs               []int
		artIDs               []int
		err                  error
		total                int64
	)
	tagIDs = make([]int, 0)
	artIDs = make([]int, 0)
	err = db.Model(&ArticleTag{}).Where("article_id = ?", id).Find(&articleTagList).Error
	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	for i := 0; i < len(articleTagList); i++ {
		tagIDs = append(tagIDs, articleTagList[i].TagID)
	}
	err = db.Model(&ArticleTag{}).Where("tag_id in ?", tagIDs).Find(&atls).Error
	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	for i := 0; i < len(atls); i++ {
		artIDs = append(artIDs, atls[i].ArticleID)
	}
	err = db.Select("article.id, title, img, created_at, updated_at, `desc`").
		Order("Created_At DESC").
		Where("id in ?", artIDs).
		Find(&articleList).Count(&total).Error
	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	return articleList, errmsg.SUCCSE, total
}
