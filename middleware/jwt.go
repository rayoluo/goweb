package middleware

import (
	"errors"
	"ginblog/utils"
	"ginblog/utils/errmsg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type JWT struct {
	JwtKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(utils.JwtKey),
	}
}

type MyClaims struct {
	// 自定义字段
	Username string `json:"username"`
	jwt.StandardClaims
}

// 定义错误
var (
	TokenInvalid error = errors.New("这不是一个token,请重新登录")
	TokenExpired error = errors.New("token过期，请重新登录")
)

// CreateToken 生成token
func (j *JWT) CreateToken(claims MyClaims) (string, error) {
	// 指定token的header和payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 对token进行签名，并存放到token的signature字段
	return token.SignedString(j.JwtKey)
}

// ParserToken 解析token
func (j *JWT) ParserToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JwtKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			}
		}
		return nil, TokenInvalid
	}

	if token != nil {
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			return claims, nil
		}
	}

	return nil, TokenInvalid
}

// JwtToken jwt中间件
// 用户未携带token字段，返回错误码errmsg.ERROR_TOKEN_EXIST = 1004
// token格式不正确，返回错误码errmsg.ERROR_TOKEN_TYPE_WRONG=1007
// token过期，返回错误码errmsg.ERROR_TOKEN_RUNTIME = 1005
// 其他错误，返回错误码errmsg.ERROR_TOKEN_WRONG=1006
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			// 用户未携带token字段，返回错误码
			code = errmsg.ERROR_TOKEN_EXIST
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		checkToken := strings.Split(tokenHeader, " ")
		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			code = errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		j := NewJWT()
		// 解析token
		claims, err := j.ParserToken(checkToken[1])
		if err != nil {
			if err == TokenExpired {
				code = errmsg.ERROR_TOKEN_RUNTIME
				c.JSON(http.StatusOK, gin.H{
					"status":  code,
					"message": "token授权已过期,请重新登录",
					"data":    nil,
				})
				c.Abort()
				return
			}
			// 其他错误
			code = errmsg.ERROR_TOKEN_WRONG
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": err.Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
