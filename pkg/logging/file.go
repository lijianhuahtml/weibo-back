package logging

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

// 函数用于获取日志文件的保存路径
func getLogFilePath() string {
	return fmt.Sprintf("%s", viper.GetString("logging.LogSavePath"))
}

// 函数用于获取完整的日志文件路径
func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s-%s.%s", viper.GetString("logging.LogSaveName"), time.Now().Format(viper.GetString("logging.TimeFormat")), viper.GetString("logging.LogFileExt"))
	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

// 函数用于打开日志文件，它接受一个文件路径作为参数，首先检查文件是否存在，
// 如果不存在则创建对应的目录结构。然后以追加、创建、写入模式打开文件，并返回文件句柄
func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}

	return handle
}

// 函数用于创建日志保存路径中的目录结构
func mkDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm) // const定义ModePerm FileMode = 0777
	if err != nil {
		panic(err)
	}
}
