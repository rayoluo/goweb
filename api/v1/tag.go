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
 * @Date: 2021/5/18 16:00
 * @Desc: 添加标签的controller层
 */

// AddTag 添加标签
func AddTag(c *gin.Context) {
	var data model.Tag
	_ = c.ShouldBindJSON(&data)
	code = model.CheckTagUsed(data.Name)
	if code == errmsg.SUCCSE {
		model.CreateTag(&data)
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// GetTagInfo 查询标签信息
func GetTagInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	data, code := model.GetTagInfo(id)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// GetTagList 查询标签列表
func GetTagList(c *gin.Context) {
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

	data, total := model.GetTagList(pageSize, pageNum)
	code = errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// GetAllTags 获取所有标签
func GetAllTags(c *gin.Context) {
	data, total := model.GetAllTags()
	code := errmsg.SUCCSE
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCSE,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// EditTag 编辑标签名
func EditTag(c *gin.Context) {
	var data model.Tag
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)
	code = model.CheckTagUsed(data.Name)
	if code == errmsg.SUCCSE {
		model.EditTag(id, &data)
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code = model.DeleteTag(id)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
