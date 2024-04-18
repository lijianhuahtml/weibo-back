package main

import (
	"weibo/pkg/mysql"
	"weibo/pkg/redis"
	"weibo/pkg/setting"
	"weibo/router"
)

func main() {
	setting.Init()
	mysql.InitMySQL()
	redis.InitRedis()

	router.Router()
}
