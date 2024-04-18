package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

var (
	Red *redis.Client
	ctx = context.Background()
)

func InitRedis() {
	fmt.Println("[redis]: init start")
	Red = redis.NewClient(&redis.Options{
		// 一旦Viper读取了配置文件，我们就可以使用Get函数来获取配置值：
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	fmt.Println("[redis]: init start")
}

// SetEmailToken 设置键值对
func SetEmailToken(key string, expiration time.Duration) error {
	return Red.Set(ctx, "email:token"+key, key, expiration).Err()
}

// EmailTokenExists 判断键是否存在
func EmailTokenExists(key string) (bool, error) {
	exists, err := Red.Exists(ctx, "email:token"+key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
