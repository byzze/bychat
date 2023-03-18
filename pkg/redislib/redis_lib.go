/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 14:18
 */

package redislib

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	client *redis.Client
)

// InitRedlisClient 初始链接
func InitRedlisClient() {
	// 初始redis链接
	client = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConns"),
	})
	// 验证是否链接成功
	pong, err := client.Ping().Result()
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info("初始化redis:", pong)
	// Output: PONG <nil>
}

// GetClient 获取redis client
func GetClient() (c *redis.Client) {
	return client
}

func CheckNilErr() {

}
