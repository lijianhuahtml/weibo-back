package mysql

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB *gorm.DB
)

func InitMySQL() {
	fmt.Println("[mysql]: init start")
	// 自定义日志模板 打印SQL语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 色彩
		},
	)
	user := viper.GetString("mysql.User")
	password := viper.GetString("mysql.Password")
	host := viper.GetString("mysql.Host")
	name := viper.GetString("mysql.Name")

	DB, _ = gorm.Open(mysql.Open(user+":"+password+"@tcp("+host+")/"+name+"?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{Logger: newLogger})

	fmt.Println("[mysql]: init finish")
}
