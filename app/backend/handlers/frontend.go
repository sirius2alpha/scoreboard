package handlers

import (
	"log"
	"net/http"
	"scoreboard/app/backend/core"

	"github.com/gin-gonic/gin"
)

func ConfigHandler(ctx *gin.Context) {
	maxKeepSeconds := core.AppConfig.GetInt("redis.max_keep_seconds")
	log.Println("max_keep_seconds: ", maxKeepSeconds)
	// 返回	Access-Control-Allow-Origin 头
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(http.StatusOK, gin.H{"max_keep_seconds": maxKeepSeconds})
}
