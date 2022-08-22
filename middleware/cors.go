package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"regexp"
)

func Cors() gin.HandlerFunc {
	// return cors.New(
	// 	cors.Config{
	// 		// AllowAllOrigins:  true,
	// 		AllowOrigins:     []string{"http://localhost:8080"},
	// 		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 		AllowHeaders:     []string{"*"},
	// 		ExposeHeaders:    []string{"Content-Length", "Authorization", "Content-Type", "Set-Cookie"},
	// 		AllowCredentials: true,
	// 		MaxAge:           12 * time.Hour,
	// 	},
	// )
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cookie", "Authorization", "X-Requested-With"}
	config.ExposeHeaders = []string{"Content-Length", "Authorization", "Content-Type", "Set-Cookie"}
	// config.AllowOrigins = []string{"http://127.0.0.1:8080", "http://localhost:8080"}
	if gin.Mode() == gin.ReleaseMode {
		// 生产环境需要配置跨域域名，否则403
		config.AllowOrigins = []string{"https://rayoluo.top", "https://admin.rayoluo.top"}
	} else {
		// 测试环境下模糊匹配本地开头的请求
		config.AllowOriginFunc = func(origin string) bool {
			if regexp.MustCompile(`^http://127\.0\.0\.1:\d+$`).MatchString(origin) {
				return true
			}
			if regexp.MustCompile(`^http://localhost:\d+$`).MatchString(origin) {
				return true
			}
			return false
		}
	}
	config.AllowCredentials = true
	return cors.New(config)
}
