package routers

import (
	"scoreboard/src/backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 定义WebSocket路由
	r.GET("/ws", controllers.HandleWebSocket)

	return r
}
