package routers

import (
	"scoreboard/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 定义WebSocket路由
	r.GET("/ws", controllers.HandleWebSocket)

	// 定义其他路由
	r.GET("/ping", controllers.Ping)

	return r
}
