package v1

import (
	"context"
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// AddArticle 添加文章-这里的逻辑需要修改，添加文章需要为文章绑定相应的标签
func AddArticle(c *gin.Context) {
	var data model.Article
	_ = c.ShouldBindJSON(&data)

	code = model.CreateArt(&data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetCateArt 查询分类下的所有文章
// func GetCateArt(c *gin.Context) {
// 	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
// 	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
// 	id, _ := strconv.Atoi(c.Param("id"))
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
// 	data, captcha, total := model.GetCateArt(id, pageSize, pageNum)
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  captcha,
// 		"data":    data,
// 		"total":   total,
// 		"message": errmsg.GetErrMsg(captcha),
// 	})
// }

// GetArtInfo 前端查询单个文章信息
func GetArtInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, code := model.GetArtInfo(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetArtInfoEdit 后端查询单个文章信息
func GetArtInfoEdit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, code := model.GetArtInfoEdit(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// // GetArt 查询文章列表
// func GetArt(c *gin.Context) {
// 	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
// 	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
// 	title := c.Query("title")
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
// 	// 如果没有指定文章标题，则调用model.GetArt
// 	if len(title) == 0 {
// 		data, captcha, total := model.GetArt(pageSize, pageNum)
// 		c.JSON(http.StatusOK, gin.H{
// 			"status":  captcha,
// 			"data":    data,
// 			"total":   total,
// 			"message": errmsg.GetErrMsg(captcha),
// 		})
// 		return
// 	}
//
// 	data, captcha, total := model.SearchArticle(title, pageSize, pageNum)
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  captcha,
// 		"data":    data,
// 		"total":   total,
// 		"message": errmsg.GetErrMsg(captcha),
// 	})
// }

// EditArt 编辑文章
func EditArt(c *gin.Context) {
	var data model.Article
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)
	// log.Println(data)
	code = model.EditArt(id, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// DeleteArt 删除文章
func DeleteArt(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code = model.DeleteArt(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetAllArt 获取所有文章
func GetAllArt(c *gin.Context) {
	data, code, total := model.GetAllArt()
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetRecArts 根据文章id获取推荐文章
func GetRecArts(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, code, total := model.GetRecArts(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetArtLikeCount 获取文章的点赞数
func GetArtLikeCount(c *gin.Context) {
	var (
		artId     int
		likeCount int64
		rdh       *model.RedisHelper
		key       string
		ctx       context.Context
		err       error
	)
	artId, _ = strconv.Atoi(c.Param("artId"))
	// 操作redis
	rdh = model.GetRedisHelper()
	key = "article:like:" + strconv.Itoa(artId)
	ctx = context.Background()
	if likeCount, err = rdh.SCard(ctx, key).Result(); err != nil {
		code = errmsg.ERROR
	} else {
		code = errmsg.SUCCSE
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    code,
		"likeCount": likeCount,
		"message":   errmsg.GetErrMsg(code),
	})
}
