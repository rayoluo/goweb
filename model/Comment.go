package model

import (
	"context"
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Comment struct {
	gorm.Model
	UserId    uint   `json:"user_id"`
	ArticleId uint   `json:"article_id"`
	Title     string `json:"article_title"`
	Username  string `json:"username"`
	Content   string `gorm:"type:varchar(500);not null;" json:"content"`
	// 评论的状态，默认显示所有评论，后台可以撤下评论, 将默认值设置为1
	Status   int8 `gorm:"type:tinyint;default:1" json:"status"`
	ReplyId  uint `json:"reply_id"`  // 回复用户id
	ParentId uint `json:"parent_id"` // 父级评论id
}

// AddComment 新增评论
func AddComment(data *Comment) int {
	err = db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	// 相应文章的评论数+1
	var art Article
	db.Model(&art).Where("id = ?", data.ArticleId).UpdateColumn("read_count", gorm.Expr("read_count + ?", 1))
	return errmsg.SUCCSE
}

// GetComment 查询单个评论
func GetComment(id int) (Comment, int) {
	var comment Comment
	err = db.Where("id = ?", id).First(&comment).Error
	if err != nil {
		return comment, errmsg.ERROR
	}
	return comment, errmsg.SUCCSE
}

// GetCommentList 后台所有获取评论列表
func GetCommentList(pageSize int, pageNum int) ([]Comment, int64, int) {
	var commentList []Comment
	var total int64
	// 获取评论总数，包含所有经过审核的、未经过审核的评论
	db.Find(&commentList).Count(&total)
	err = db.Model(&commentList).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("Created_At DESC").Select("comment.id, article.title,user_id,article_id, user.username, comment.content, comment.status,comment.created_at,comment.deleted_at").Joins("LEFT JOIN article ON comment.article_id = article.id").Joins("LEFT JOIN user ON comment.user_id = user.id").Scan(&commentList).Error
	if err != nil {
		return commentList, 0, errmsg.ERROR
	}
	return commentList, total, errmsg.SUCCSE
}

// GetCommentCount 获取评论数量
func GetCommentCount(id int) int64 {
	var commentList []Comment
	var total int64
	db.Where("article_id = ?", id).Where("status = ?", 1).Find(&commentList).Count(&total)
	return total
}

type Record struct {
	ID             uint          `json:"id"`
	UserId         uint          `json:"userId"`
	CommentContent string        `json:"commentContent"`
	CreateTime     time.Time     `json:"createTime"`
	Avatar         string        `json:"avatar"`
	Nickname       string        `json:"nickname"`
	Website        string        `json:"website"`
	ReplyCount     int64         `json:"replyCount"`
	ReplyDTOList   []ReplyRecord `json:"replyDTOList"`
	LikeCount      int64         `json:"likeCount"`
}
type ReplyRecord struct {
	ID             uint      `json:"id"`
	UserId         uint      `json:"userId"`
	CommentContent string    `json:"commentContent"`
	CreateTime     time.Time `json:"createTime"`
	Avatar         string    `json:"avatar"`
	Nickname       string    `json:"nickname"`
	Website        string    `json:"website"`
	LikeCount      int64     `json:"likeCount"`
	ParentId       uint      `json:"parentId"`
	ReplyId        uint      `json:"replyId"`
	ReplyNickname  string    `json:"replyNickname"`
	ReplyWebSite   string    `json:"replyWebSite"`
}

func GetCommentReplies(commentId, pageSize, pageNum int) ([]ReplyRecord, int) {
	var (
		replyRecList  []ReplyRecord
		childComments []Comment
		j             int
		userInfo      User
		err           error
	)
	// 根据父级评论ID获取所有的回复
	// 找到所有parentId为当前评论id的所有评论
	err = db.Model(&Comment{}).Select("id, user_id, created_at, content, parent_id, reply_id").
		Limit(pageSize).Offset((pageNum-1)*pageSize).
		Where("parent_id = ?", commentId).
		Where("status = ?", 1).Find(&childComments).Error
	if err != nil {
		return nil, errmsg.ERROR
	}

	for j = 0; j < len(childComments); j++ {
		var replyRec ReplyRecord
		replyRec.ID = childComments[j].ID
		replyRec.UserId = childComments[j].UserId
		replyRec.CreateTime = childComments[j].CreatedAt
		replyRec.CommentContent = childComments[j].Content
		replyRec.ParentId = childComments[j].ParentId
		replyRec.ReplyId = childComments[j].ReplyId

		db.Model(&User{}).Select("avatar", "nickname", "website").
			Where("id = ?", replyRec.UserId).First(&userInfo)
		replyRec.Avatar = userInfo.Avatar
		replyRec.Nickname = userInfo.Nickname
		replyRec.Website = userInfo.Website
		db.Model(&User{}).Select("nickname", "website").
			Where("id = ?", replyRec.ReplyId).First(&userInfo)
		replyRec.ReplyNickname = userInfo.Nickname
		replyRec.ReplyWebSite = userInfo.Website

		// todo 点赞功能统计
		replyRec.LikeCount = 0

		replyRecList = append(replyRecList, replyRec)
	}

	return replyRecList, errmsg.SUCCSE
}

func GetRecordList(id int, pageSize int, pageNum int) ([]Record, int64, int) {
	var (
		commentList, childComments []Comment
		total                      int64
		code, i, j                 int
		recordList                 []Record
		userInfo                   User
	)
	// 获取顶级评论
	commentList, total, code = GetCommentListFront(id, pageSize, pageNum)
	// log.Println(commentList)
	for i = 0; i < len(commentList); i++ {
		var rec Record
		rec.ID = commentList[i].ID
		rec.UserId = commentList[i].UserId
		rec.CreateTime = commentList[i].CreatedAt
		rec.CommentContent = commentList[i].Content
		db.Model(&User{}).Select("avatar", "nickname", "website").
			Where("id = ?", rec.UserId).First(&userInfo)
		rec.Avatar = userInfo.Avatar
		rec.Nickname = userInfo.Nickname
		rec.Website = userInfo.Website

		// 找到所有parentId为当前评论id的所有评论
		db.Model(&Comment{}).Select("id, user_id, created_at, content, parent_id, reply_id").
			Where("parent_id = ?", rec.ID).
			Where("status = ?", 1).Find(&childComments).Count(&rec.ReplyCount)
		// 只展示前三条
		for j = 0; j < len(childComments) && j < 3; j++ {
			var replyRec ReplyRecord
			replyRec.ID = childComments[j].ID
			replyRec.UserId = childComments[j].UserId
			replyRec.CreateTime = childComments[j].CreatedAt
			replyRec.CommentContent = childComments[j].Content
			replyRec.ParentId = childComments[j].ParentId
			replyRec.ReplyId = childComments[j].ReplyId

			db.Model(&User{}).Select("avatar", "nickname", "website").
				Where("id = ?", replyRec.UserId).First(&userInfo)
			replyRec.Avatar = userInfo.Avatar
			replyRec.Nickname = userInfo.Nickname
			replyRec.Website = userInfo.Website
			db.Model(&User{}).Select("nickname", "website").
				Where("id = ?", replyRec.ReplyId).First(&userInfo)
			replyRec.ReplyNickname = userInfo.Nickname
			replyRec.ReplyWebSite = userInfo.Website

			// todo 点赞功能统计
			replyRec.LikeCount = GetCommentLikeCount(int(replyRec.ID))

			rec.ReplyDTOList = append(rec.ReplyDTOList, replyRec)
		}

		// todo 点赞功能统计 从redis中获取
		rec.LikeCount = GetCommentLikeCount(int(rec.ID))

		recordList = append(recordList, rec)
	}
	return recordList, total, code
}

// GetCommentLikeCount 获取评论的点赞数
func GetCommentLikeCount(commentId int) int64 {
	var (
		rdh       *RedisHelper
		key       string
		ctx       context.Context
		likeCount int64
	)
	rdh = GetRedisHelper()
	key = "comment:like:" + strconv.Itoa(commentId)
	ctx = context.Background()
	likeCount, _ = rdh.SCard(ctx, key).Result()
	return likeCount
}

// GetCommentListFront 获取所有的父级评论
func GetCommentListFront(id int, pageSize int, pageNum int) ([]Comment, int64, int) {
	var commentList []Comment
	var total int64
	db.Find(&Comment{}).Where("article_id = ?", id).Where("parent_id = ?", 0).Where("status = ?", 1).Count(&total)
	// err = db.Model(&Comment{}).Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Select("comment.id, comment.user_id, article_id, user.username, comment.content, comment.status,comment.created_at,comment.deleted_at").Joins("LEFT JOIN article ON comment.article_id = article.id").Joins("LEFT JOIN user ON comment.user_id = user.id").Where("article_id = ?",
	// 	id).Where("status = ?", 1).Scan(&commentList).Error
	err = db.Model(&Comment{}).Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").
		Select("id, user_id, created_at, content").
		Where("article_id = ?", id).
		Where("parent_id = ?", 0).
		Where("status = ?", 1).Find(&commentList).Error
	if err != nil {
		return commentList, 0, errmsg.ERROR
	}
	return commentList, total, errmsg.SUCCSE
}

// 编辑评论（暂不允许编辑评论）

// DeleteComment 删除评论
func DeleteComment(id uint) int {
	var comment Comment
	err = db.Where("id = ?", id).Delete(&comment).Error
	if err != nil {
		return errmsg.ERROR
	}
	// 如果当前评论是通过审核被删除，则评论数应该减1
	if comment.Status == 1 {
		db.Model(&Article{}).Where("id = ?", comment.ArticleId).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	}
	return errmsg.SUCCSE
}

// CheckComment 通过评论
func CheckComment(id int, data *Comment) int {
	var comment Comment
	var res Comment
	var article Article
	var maps = make(map[string]interface{})
	maps["status"] = data.Status // =1

	// 获取审核状态更新完成的评论之后需要将该评论所在的文章的评论数加1
	err = db.Model(&comment).Where("id = ?", id).Updates(maps).First(&res).Error
	db.Model(&article).Where("id = ?", res.ArticleId).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// UncheckComment 撤下评论
func UncheckComment(id int, data *Comment) int {
	var comment Comment
	var res Comment
	var article Article
	var maps = make(map[string]interface{})
	maps["status"] = data.Status

	err = db.Model(&comment).Where("id = ?", id).Updates(maps).First(&res).Error
	db.Model(&article).Where("id = ?", res.ArticleId).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
