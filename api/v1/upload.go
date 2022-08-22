package v1

import (
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpLoad(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")

	fileSize := fileHeader.Size

	url, code := model.UpLoadFile(file, fileSize)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"url":     url,
	})
}

// // UpLoadFront 前端文件上传
// func UpLoadFront(c *gin.Context) {
// 	// log.Println(c.Request.Header)
// 	// 获取上传的文件
// 	var (
// 		code int
// 		file         multipart.File
// 		fileHeader   *multipart.FileHeader
// 		url          string
// 		fileSize     int64
// 		err          error
// 	)
//
// 	if file, fileHeader, err = c.Request.FormFile("file"); err != nil {
// 		code = errmsg.ERROR
// 		c.JSON(http.StatusOK, gin.H{
// 			"status":  code,
// 			"message": errmsg.GetErrMsg(code),
// 			"url":     "",
// 		})
// 		return
// 	}
// 	fileSize = fileHeader.Size
//
// 	// if code = CheckSessionInfo(c, userId); code != errmsg.SUCCSE {
// 	// 	c.JSON(http.StatusOK, gin.H{
// 	// 		"status":  code,
// 	// 		"message": errmsg.GetErrMsg(code),
// 	// 		"url":     "",
// 	// 	})
// 	// 	return
// 	// }
//
// 	url, code = model.UpLoadFile(file, fileSize)
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  code,
// 		"message": errmsg.GetErrMsg(code),
// 		"url":     url,
// 	})
// }
