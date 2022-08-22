package model

import (
	"context"
	"ginblog/utils"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"time"
)

/**
 * @Author: rayoluo
 * @Author: hustly123@gmail.com
 * @Date: 2021/6/3 16:46
 * @Desc: redis模块，连接redis数据库，获取redis实例
 */

type RedisHelper struct {
	*redis.Client
}

var redisHelper *RedisHelper

var redisOnce sync.Once

func GetRedisHelper() *RedisHelper {
	return redisHelper
}

func newRedisHelper() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         utils.RedisAddr + ":" + utils.RedisPort,
		Password:     "",
		DB:           utils.RedisDB,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	redisOnce.Do(func() {
		rdh := new(RedisHelper)
		rdh.Client = rdb
		redisHelper = rdh
	})
	return rdb
}

func InitRedis() {
	var ctx = context.Background()
	rdb := newRedisHelper()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal(err.Error())
		return
	}
}
