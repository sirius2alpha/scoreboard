package main

import (
	"fmt"
	"log"
	"net/http"
	"scoreboard/src/backend/routers"

	"github.com/spf13/viper"
)

func main() {
	// 从环境变量中读取后端监听的端口
	viper.SetConfigName("config") // 配置文件名称（无扩展名）
	viper.SetConfigType("toml")   // 或viper.SetConfigType("INI")
	viper.AddConfigPath(".")      // 配置文件路径

	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// 直接通过viper获取配置
	backendPort := viper.GetString("backend_listen_port")

	fmt.Println("端口是", backendPort)

	// 路由初始化
	router := routers.SetupRouter()

	// 启动服务器
	log.Fatal(http.ListenAndServe(":"+backendPort, router))
}
