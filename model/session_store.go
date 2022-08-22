package model

import (
	"ginblog/utils"
	"github.com/gin-contrib/sessions/redis"
)

/**
 * @Author: rayoluo
 * @Author: hustly123@gmail.com
 * @Date: 2021/6/5 10:48
 * @Desc: 返回session store
 */

var store redis.Store

func GetSessionStore() redis.Store {
	return store
}

func InitSessionStore() {
	var addr = utils.StoreAddr + ":" + utils.StorePort
	// 初始化基于redis的存储引擎
	// 参数说明：
	//    第1个参数 - redis最大的空闲连接数
	//    第2个参数 - 数通信协议tcp或者udp
	//    第3个参数 - redis地址, 格式，host:port
	//    第4个参数 - redis密码
	//    第5个参数 - session加密密钥
	store, _ = redis.NewStore(utils.IdleConnSize, utils.SessionNetwork, addr, "", []byte(utils.StoreSecret))
}
