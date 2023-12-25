package main

import (
    "github.com/gin-gonic/gin"
    "github.com/sirius2alpha/scoreboard/routers"
)

func main() {
    r := gin.Default()
    routers.SetupRouter(r)
    r.Run()
}
