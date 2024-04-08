package routers

import (
	"scoreboard/app/backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 定义WebSocket路由
	r.GET("/ws", handlers.HandleWebSocket)
	r.GET("/config", handlers.ConfigHandler)

	return r
}
