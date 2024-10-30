package cache

import (
	"Take_Out/config"
	"fmt"
	"github.com/go-redis/redis"
)

var (
	RedisClient *redis.Client
)

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			config.Config.Redis.RedisHost,
			config.Config.Redis.RedisPort,
		),
		Password: config.Config.Redis.RedisPassword,
		//DB:       config.Config.Redis.RedisDbName,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("%s\n", err))
	}
	RedisClient = client
}
