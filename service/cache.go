package service

import (
	"fcbox/config"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var redisClient *redis.Client

const (
	CachePrefix = "FCBOX_CACHE_PREFIX_"
	Yes         = 1
	No          = 0
)

func init() {
	yaml := config.ReadYaml().Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     yaml.Address + ":" + yaml.Port,
		Password: yaml.Password, // no password set
		DB:       0,             // use default DB
	})
}

// IsSend 是否需要发送短信
func IsSend(code string) bool {
	key := buildKey(code)
	exists := redisClient.Exists(redisClient.Context(), key)
	if exists.Err() != nil {
		// 存在错误 直接宕机
		ErrorBark("redis exist error")
		panic(exists.Err())
	}
	result, err := exists.Result()
	if err != nil {
		ErrorBark("redis exist result error")
		panic(err)
	}
	// 不存在 即需要发送
	return result == No
}

func buildKey(code string) string {
	key := CachePrefix + code
	log.Printf("build cache key:%s\n", key)
	return key
}

func Set(code string) {
	key := buildKey(code)
	frequency := config.ReadYaml().Notify.Frequency
	expire := time.Duration(frequency) * time.Hour
	redisClient.SetNX(redisClient.Context(), key, "1", expire)
}
