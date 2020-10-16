package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	Pool      *redis.Pool
	redisHost = "ip:port"
	redisPass = "password"
)

// redis conn pool
func NewRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			// 1、打开链接
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			// 2、访问认证
			if _, err := c.Do("AUTH", redisPass); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func init() {
	Pool = NewRedisPool()
}

func RedisPool() *redis.Pool {
	return Pool
}
