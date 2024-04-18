package setting

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitConfig() {
	fmt.Println("[setting]: init start")
	viper.SetConfigName("app")      // 配置文件名（不带后缀）
	viper.AddConfigPath("./config") // 查找配置文件所在的路径，多次调用以添加多个搜索路径
	err := viper.ReadInConfig()     // 查找并读取配置文件
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("[setting]: init finish")
}

func init() {
	InitConfig()
}

func Init() {}
