package v1

import (
	"encoding/json"
	"ginblog/middleware"
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Login 后台登陆
func Login(c *gin.Context) {
	var formData model.User
	_ = c.ShouldBindJSON(&formData)
	var token string
	var code int

	formData, code = model.CheckLogin(formData.Username, formData.Password)

	// 如果成功登录就为其分配一个token
	if code == errmsg.SUCCSE {
		// 设置token并返回
		setToken(c, formData)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    formData.Username,
			"id":      formData.ID,
			"message": errmsg.GetErrMsg(code),
			"token":   token,
		})
	}

}

// LoginFront 前台登录
func LoginFront(c *gin.Context) {
	var formData model.User
	_ = c.ShouldBindJSON(&formData)

	var code int
	var data *model.UserWithLikeSet

	// 首先验证登录用户名和密码是否正确
	formData, code = model.CheckLoginFront(formData.Username, formData.Password)
	// 如果验证通过则设置session信息
	if code == errmsg.SUCCSE {
		sessionInfo, _ := json.Marshal(formData)
		session := sessions.Default(c)
		// 设置session信息，即用户名和用户信息键值对
		// 这里没有save...
		// log.Println("设置session和持久化")
		session.Set("user_"+string(formData.ID), sessionInfo)
		session.Save()
		if err := session.Save(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR,
				"data":    "",
				"message": errmsg.GetErrMsg(errmsg.ERROR),
			})
			return
		}
		// 获取用户点赞集合
		data = model.GetUserLikeSet(&formData)

		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    *data,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    nil,
		"message": errmsg.GetErrMsg(code),
	})
}

// Logout 前端退出登录
func Logout(c *gin.Context) {
	var (
		sess sessions.Session
		err  error
		code int
	)
	sess = sessions.Default(c)
	sess.Clear()
	// 删除当前session
	// sess.Options(sessions.Options{MaxAge: -1})
	if err = sess.Save(); err != nil {
		code = errmsg.ERROR_SESSION_REMOVE
	} else {
		code = errmsg.SUCCSE
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// token生成函数
func setToken(c *gin.Context, user model.User) {
	j := middleware.NewJWT()
	claims := middleware.MyClaims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 100,  // 生效时间
			ExpiresAt: time.Now().Unix() + 7200, // 过期时间
			Issuer:    "YuGoBlog",
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCSE,
		"data":    user.Username,
		"id":      user.ID,
		"message": errmsg.GetErrMsg(errmsg.SUCCSE),
		"token":   token,
	})
	return
}
