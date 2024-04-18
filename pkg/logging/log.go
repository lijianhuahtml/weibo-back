package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// Level 定义日志级别常量
type Level int

// 日志文件句柄
var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"} // 日志级别标志
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func init() {
	fmt.Println("[logging]: init start")
	// 获取日志文件完整路径
	filePath := getLogFileFullPath()
	// 打开日志文件
	F = openLogFile(filePath)
	// 创建 Logger 实例, log.LstdFlags：日志记录的格式属性之一, 其余的选项如下
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
	fmt.Println("[logging]: init finish")
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}

// 设置日志前缀
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
