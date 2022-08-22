package captcha

import (
	"math/rand"
	"time"
)

/**
 * @Author: rayoluo
 * @Author: hustly123@gmail.com
 * @Date: 2021/6/2 23:49
 * @Desc: 生成随机验证码
 */

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// 生成指定长度的验证码
func CreateCaptcha(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
