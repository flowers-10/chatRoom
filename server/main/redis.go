package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// 定一个全局pool
var pool *redis.Pool

func initPool(address string, maxId, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxId,       //最大空闲链接数
		MaxActive:   maxActive,   //和数据库最大链接数，0没限制
		IdleTimeout: idleTimeout, //最大空闲时间
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}
