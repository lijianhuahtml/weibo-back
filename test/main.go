package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"weibo/models"
)

func main() {
	db, _ := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/wb?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	err := db.AutoMigrate(&models.Account{})
	if err != nil {
		return
	}
}
