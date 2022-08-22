package v1

import (
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/**
 * @Author: rayoluo
 * @Author: hustly123@gmail.com
 * @Date: 2021/5/23 1:23
 * @Desc: 文章和标签对应处理函数
 */

// AddArtTags 添加文章和标签的对应关系 文章id在url中指定
func AddArtTags(c *gin.Context) {
	articleID, _ := strconv.Atoi(c.Param("id"))
	var tagIDs []int
	c.ShouldBindJSON(&tagIDs)
	code := model.AddArtTags(articleID, tagIDs)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// DeleteByArtID 根据文章ID删除对应关系
func DeleteByArtID(c *gin.Context) {
	articleID, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteByArtID(articleID)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// DeleteByTagID 根据标签ID删除对应关系
func DeleteByTagID(c *gin.Context) {
	tagID, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteByTagID(tagID)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetTagsByArtID 根据文章ID获取标签列表
func GetTagsByArtID(c *gin.Context) {
	articleID, _ := strconv.Atoi(c.Param("id"))
	data, code, total := model.GetArtTags(articleID)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetArtsByTagID 根据标签ID获取文章列表
func GetArtsByTagID(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	tagID, _ := strconv.Atoi(c.Param("id"))

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, code, total := model.GetArtsByTag(tagID, pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetArtsWithTags 获取文章标签列表
func GetArtsWithTags(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	title := c.Query("title")

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}
	// 如果没有指定文章标题，则调用model.GetArt
	if len(title) == 0 {
		data, code, total := model.GetArtTagsList(pageSize, pageNum)
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	data, code, total := model.SearchArtTags(title, pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetCateArtTags 查询分类下的所有文章+标签
func GetCateArtTags(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	id, _ := strconv.Atoi(c.Param("id"))

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, code, total := model.GetCateArtTags(id, pageSize, pageNum)
	// log.Println(data)
	// str, _ := json.Marshal(data)
	// log.Println(string(str))

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// 思考在增加标签功能之后哪里需要修改
// --对于新增文章，在数据库post文章之后可以获取到文章ID，根据文章
// ID和选中的标签，在标签-文章表中增加相应数据
// --对于编辑文章，在数据库edit文章之后，首先根据文章的ID删除
// 当前文章的所有标签，然后为当前文章添加标签 执行的逻辑操作和新增文章相同
// --对于获取文章，或许需要一个接口能够获取所有文章和对应的标签，前端需要
// --对于删除文章，需要根据文章的ID删除对应的标签

// 删除分类之后，文章并没有被删除，在删除标签之后，文章
// 也不应该被删除，但是应该删除该标签的所有对应关系
// 这些功能不应该由后端实现，阶段直接调相应的接口实现操作就可以了
