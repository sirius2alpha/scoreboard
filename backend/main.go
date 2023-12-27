package main

import (
	"log"
	"net/http"
	"scoreboard/routers"
)

func main() {

	// 路由初始化
	router := routers.SetupRouter()

	// 启动服务器
	log.Fatal(http.ListenAndServe(":8080", router))
}
