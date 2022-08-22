package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"ginblog/model"
	"ginblog/utils"
	"ginblog/utils/captcha"
	"ginblog/utils/email"
	"ginblog/utils/errmsg"
	"ginblog/utils/randstr"
	"ginblog/utils/validator"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"time"
)

var code int

// AddUser 添加用户
func AddUser(c *gin.Context) {
	var (
		data                           model.User
		msg, inputCode, key, rightCode string
		validCode, retCode             int
		redisHelper                    *model.RedisHelper
		ctx                            context.Context
		err                            error
		rawData                        []byte
		m                              map[string]interface{}
	)

	rawData, _ = c.GetRawData()
	// 反序列化
	_ = json.Unmarshal(rawData, &m)
	data.Username = m["username"].(string)
	data.Password = m["password"].(string)
	data.Role = 2
	data.Avatar = utils.Avatar
	// 为用户生成初始随机用户名
	data.Nickname = "用户" + randstr.GetRandomString(12)

	inputCode = m["code"].(string)

	// 验证验证码是否正确
	redisHelper = model.GetRedisHelper()
	ctx = context.Background()
	key = "user:info:mail:" + data.Username
	rightCode, err = redisHelper.Get(ctx, key).Result()
	// log.Println(rightCode)
	// log.Println(data.Username)
	// log.Println(data.Password)
	// log.Println(inputCode)
	if err != nil {
		if err == redis.Nil { // 验证码失效
			retCode = errmsg.ERROR_CAPTCHA_RUNTIME
		} else {
			retCode = errmsg.ERROR
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  retCode,
			"message": errmsg.GetErrMsg(retCode),
		})
		// c.Abort()
		return
	} else { // 获取到了redis中的验证码
		if inputCode != rightCode {
			retCode = errmsg.ERROR_CAPTCHA_INCORRECT
			c.JSON(http.StatusOK, gin.H{
				"status":  retCode,
				"message": errmsg.GetErrMsg(retCode),
			})
			return
		}
	}

	msg, validCode = validator.Validate(&data)
	if validCode != errmsg.SUCCSE {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  validCode,
				"message": msg,
			},
		)
		// c.Abort()
		return
	}

	code = model.CheckUser(data.Username)
	if code == errmsg.SUCCSE {
		model.CreateUser(&data)
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// ResetPassword 前端用户重置密码
func ResetPassword(c *gin.Context) {
	// 获取前端请求的参数
	var (
		username, password, code, key, rightCode string
		redisHelper                              *model.RedisHelper
		ctx                                      context.Context
		err                                      error
		retCode                                  int
		rawData                                  []byte
		m                                        map[string]interface{}
	)

	rawData, _ = c.GetRawData()
	// 反序列化
	_ = json.Unmarshal(rawData, &m)
	username = m["username"].(string)
	password = m["password"].(string)
	code = m["code"].(string)

	// 验证用户提交验证码是否正确
	redisHelper = model.GetRedisHelper()
	ctx = context.Background()
	key = "user:info:mail:" + username
	rightCode, err = redisHelper.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil { // 验证码失效
			retCode = errmsg.ERROR_CAPTCHA_RUNTIME
		} else {
			retCode = errmsg.ERROR
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  retCode,
			"message": errmsg.GetErrMsg(retCode),
		})
		// c.Abort()
		return
	} else { // 获取到了redis中的验证码
		if code != rightCode {
			retCode = errmsg.ERROR_CAPTCHA_INCORRECT
			c.JSON(http.StatusOK, gin.H{
				"status":  retCode,
				"message": errmsg.GetErrMsg(retCode),
			})
			return
		}
	}
	// 修改用户密码的逻辑
	retCode = model.ResetPasswordFront(username, password)
	c.JSON(http.StatusOK, gin.H{
		"status":  retCode,
		"message": errmsg.GetErrMsg(retCode),
	})
}

// GetUserInfo 查询单个用户
func GetUserInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var maps = make(map[string]interface{})
	data, code := model.GetUser(id)
	maps["username"] = data.Username
	maps["role"] = data.Role
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    maps,
			"total":   1,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// GetUsers 查询用户列表
func GetUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	username := c.Query("username")

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, total := model.GetUsers(username, pageSize, pageNum)

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

// EditUser 编辑用户
func EditUser(c *gin.Context) {
	var data model.User
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code = model.CheckUpUser(id, data.Username)
	if code == errmsg.SUCCSE {
		model.EditUser(id, &data)
	}
	if code == errmsg.ERROR_USERNAME_USED {
		c.Abort()
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// 修改密码
func ChangeUserPassword(c *gin.Context) {
	var data model.User
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code = model.ChangePassword(id, &data)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// 删除用户
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code = model.DeleteUser(id)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// GenCaptcha 请求发送验证码到指定的邮箱
func GenCaptcha(c *gin.Context) {
	var (
		username, code, key, subject, content string
		redisHelper                           *model.RedisHelper
		ctx                                   context.Context
		err                                   error
	)
	// 获取用户名，即邮箱
	username = c.Query("username")
	// 生成随机验证码
	code = captcha.CreateCaptcha(6)
	// 将用户名和对应的邮箱验证码存入到redis中
	redisHelper = model.GetRedisHelper()
	ctx = context.Background()
	key = "user:info:mail:" + username
	// 设置有效时间为十分钟
	redisHelper.Set(ctx, key, code, 10*time.Minute)
	// 发送邮件
	subject = "博客验证码"
	content = fmt.Sprintf("您的验证码为%s, 有效时间为%d分钟！", code, 10)
	err = email.SendEmail(subject, content, username)
	// 返回
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.SUCCSE,
			"message": errmsg.GetErrMsg(errmsg.SUCCSE),
		})
	}
}

// UpdateUserInfo 前端用户中心更新除了ID之外的个人信息：包括avatar, nickname, intro, website
func UpdateUserInfo(c *gin.Context) {
	var (
		userId, code int
		rawData      []byte
		m            map[string]interface{}
	)
	userId, _ = strconv.Atoi(c.Param("id"))
	rawData, _ = c.GetRawData()

	_ = json.Unmarshal(rawData, &m)

	// 检查用户ID是否一致或者session是否过期
	if code, _ = CheckSessionInfo(c, userId); code != errmsg.SUCCSE {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	code = model.UpdateUserInfoFront(userId, m)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// SaveArticleLike 用户给指定的文章点赞
func SaveArticleLike(c *gin.Context) {
	var (
		artId, userId int
		rawData       []byte
		m             map[string]interface{}
		rdh           *model.RedisHelper
		key           string
		ctx           context.Context
		isMember      bool
		err           error
	)
	// 获取文章的id
	artId, _ = strconv.Atoi(c.Param("artId"))
	rawData, _ = c.GetRawData()
	json.Unmarshal(rawData, &m)
	userId = int(m["userId"].(float64))
	// 验证用户session信息
	if code, _ = CheckSessionInfo(c, userId); code != errmsg.SUCCSE {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	// 操作redis
	rdh = model.GetRedisHelper()
	key = "article:like:" + strconv.Itoa(artId)
	ctx = context.Background()
	if isMember, err = rdh.SIsMember(ctx, key, userId).Result(); err != nil {
		code = errmsg.ERROR
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	if !isMember { // 用户不在集合中
		if _, err = rdh.SAdd(ctx, key, userId).Result(); err != nil {
			code = errmsg.ERROR
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			return
		}
	} else { // 用户在集合中
		if _, err = rdh.SRem(ctx, key, userId).Result(); err != nil {
			code = errmsg.ERROR
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			return
		}
	}
	code = errmsg.SUCCSE
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// SaveCommentLike 用户给指定的评论点赞
func SaveCommentLike(c *gin.Context) {
	var (
		commentId, userId int
		rawData           []byte
		m                 map[string]interface{}
		rdh               *model.RedisHelper
		key               string
		ctx               context.Context
		isMember          bool
		err               error
	)
	// 获取文章的id
	commentId, _ = strconv.Atoi(c.Param("commentId"))
	rawData, _ = c.GetRawData()
	json.Unmarshal(rawData, &m)
	userId = int(m["userId"].(float64))
	// 验证用户session信息
	if code, _ = CheckSessionInfo(c, userId); code != errmsg.SUCCSE {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	// 操作redis
	rdh = model.GetRedisHelper()
	key = "comment:like:" + strconv.Itoa(commentId)
	ctx = context.Background()
	if isMember, err = rdh.SIsMember(ctx, key, userId).Result(); err != nil {
		code = errmsg.ERROR
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	if !isMember { // 用户不在集合中
		if _, err = rdh.SAdd(ctx, key, userId).Result(); err != nil {
			code = errmsg.ERROR
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			return
		}
	} else { // 用户在集合中
		if _, err = rdh.SRem(ctx, key, userId).Result(); err != nil {
			code = errmsg.ERROR
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			return
		}
	}
	code = errmsg.SUCCSE
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// CheckSessionInfo 检查前端用户的session信息
func CheckSessionInfo(c *gin.Context, userId int) (int, *model.User) {
	// 从session中获取用户信息
	session := sessions.Default(c)
	// log.Println(userId)
	v := session.Get("user_" + string(userId))
	if v == nil {
		code = errmsg.ERROR_SESSION_RUNTIME
		return code, nil
	}
	var userInfo model.User
	if err := json.Unmarshal(v.([]byte), &userInfo); err != nil {
		code = errmsg.ERROR
		return code, nil
	}
	// log.Println(userInfo)
	// 验证服务器存的ID和请求ID是否一致，防止用户人为地修改请求更新了别的用户地信息
	if userId != int(userInfo.ID) {
		code = errmsg.ERROR_SESSION_KEY
		return code, nil
	}
	return errmsg.SUCCSE, &userInfo
}
