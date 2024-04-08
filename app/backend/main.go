package main

import (
	"fmt"
	"log"
	"net/http"
	"scoreboard/app/backend/core"
	"scoreboard/app/backend/routers"
)

func main() {

	backendPort := core.AppConfig.GetInt("server.backend_listen_port")

	log.Printf("Backend service statrted, listening on port：%d", backendPort)

	// 路由初始化
	router := routers.SetupRouter()

	// 启动服务器
	log.Fatal(http.ListenAndServe(":"+fmt.Sprintf("%d", backendPort), router))
}
