package randstr

import (
	"crypto/rand"
	"fmt"
)

/**
 * @Author: rayoluo
 * @Author: hustly123@gmail.com
 * @Date: 2021/6/7 9:42
 * @Desc: 生成随机用户名字符串
 */

func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}
