package main

import (
	"ginblog/model"
	"ginblog/routes"
)

func main() {
	// 引用数据库
	model.InitDb()
	// 初始化redis
	model.InitRedis()
	// 初始化redis store
	model.InitSessionStore()
	// 引入路由组件
	routes.InitRouter()
}
