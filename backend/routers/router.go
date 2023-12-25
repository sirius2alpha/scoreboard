package routers

import (
    "../controllers"
    "github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
    r.GET("/ping", controllers.Ping)
}
