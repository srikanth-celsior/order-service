package utils

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() *redis.Client {
	secrets, _ := GetSecrets([]string{"REDIS_USERNAME", "REDIS_PASSWORD"}, os.Getenv("PUBSUB_PROJECT_ID"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-16214.c15.us-east-1-4.ec2.redns.redis-cloud.com:16214",
		Username: secrets["REDIS_USERNAME"],
		Password: secrets["REDIS_PASSWORD"],
		DB:       0,
	})

	return rdb
}

func GetRedisClient() *redis.Client {
	if RedisClient == nil {
		RedisClient = InitRedis()
	}
	return RedisClient
}
