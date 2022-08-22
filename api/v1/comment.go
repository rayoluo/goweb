package v1

import (
	"encoding/json"
	"fmt"
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// AddComment 前端新增评论接口
// 功能描述：添加评论要求用户登录，用户登录成功之后会将用户的id,username,intro,website,avatar,nickname存到redis
// AddComment首先从session里面拿到userid,articleId和content是前端提交的
func AddComment(c *gin.Context) {
	var (
		userId, code, artId int
		userInfo            *model.User
		rawData             []byte
		m                   map[string]interface{}
		ok                  bool
		val                 interface{}
	)
	userId, _ = strconv.Atoi(c.Param("id"))
	rawData, _ = c.GetRawData()
	// 反序列化
	_ = json.Unmarshal(rawData, &m)

	// 检查用户ID是否一致或者session是否过期
	if code, userInfo = CheckSessionInfo(c, userId); code != errmsg.SUCCSE {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    "",
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	// 构造评论数据，即Comment结构体
	var data model.Comment
	data.UserId = uint(userId)
	artId, _ = strconv.Atoi(m["articleId"].(string))
	data.ArticleId = uint(artId)
	data.Content = m["commentContent"].(string)
	data.Username = userInfo.Username
	data.Title = fmt.Sprintf("来自%s的评论", userInfo.Nickname)

	if val, ok = m["replyId"]; ok {
		// todo 添加回复的邮件通知功能
		data.ReplyId = uint(val.(float64))
	}
	if val, ok = m["parentId"]; ok {
		data.ParentId = uint(val.(float64))
	}

	code = model.AddComment(&data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetComment 获取单个评论信息
func GetComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, code := model.GetComment(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code = model.DeleteComment(uint(id))
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetCommentCount 获取评论数量 前端接口 id为文章id
func GetCommentCount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	total := model.GetCommentCount(id)
	c.JSON(http.StatusOK, gin.H{
		"total": total,
	})
}

// GetCommentList 后台查询评论列表
func GetCommentList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, total, code := model.GetCommentList(pageSize, pageNum)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})

}

// GetCommentListFront 展示页面显示评论列表, 前端参数包括当前页面current和articleId
func GetCommentListFront(c *gin.Context) {
	var (
		pageNum, pageSize, artId, code int
		total                          int64
		recordList                     []model.Record
	)
	pageNum, _ = strconv.Atoi(c.Query("current"))
	artId, _ = strconv.Atoi(c.Query("articleId"))
	// 显示10条评论
	pageSize = 10

	recordList, total, code = model.GetRecordList(artId, pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    recordList,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetCommentReplies 前端获取父级评论下面的回复，按照分页进行显示
func GetCommentReplies(c *gin.Context) {
	var (
		commentId, pageSize, pageNum, code int
		replyRecs                          []model.ReplyRecord
	)
	commentId, _ = strconv.Atoi(c.Param("commentId"))
	pageNum, _ = strconv.Atoi(c.Query("current"))
	// 回复分页中每页显示5条回复
	pageSize = 5
	replyRecs, code = model.GetCommentReplies(commentId, pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    replyRecs,
		"message": errmsg.GetErrMsg(code),
	})
}

// func GetCommentListFront(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("id"))
// 	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
// 	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
//
// 	switch {
// 	case pageSize >= 100:
// 		pageSize = 100
// 	case pageSize <= 0:
// 		pageSize = 10
// 	}
//
// 	if pageNum == 0 {
// 		pageNum = 1
// 	}
//
// 	data, total, code := model.GetCommentListFront(id, pageSize, pageNum)
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  code,
// 		"data":    data,
// 		"total":   total,
// 		"message": errmsg.GetErrMsg(code),
// 	})
//
// }

// CheckComment 通过审核
func CheckComment(c *gin.Context) {
	var data model.Comment
	_ = c.ShouldBindJSON(&data)
	id, _ := strconv.Atoi(c.Param("id"))

	code = model.CheckComment(id, &data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 撤下评论审核
func UnCheckcomment(c *gin.Context) {
	var data model.Comment
	_ = c.ShouldBindJSON(&data)
	id, _ := strconv.Atoi(c.Param("id"))

	code = model.UncheckComment(id, &data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
