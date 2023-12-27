package controllers

import (
    "github.com/gin-gonic/gin"
)

// Ping 是一个处理器函数，用于处理来自客户端的 ping 请求。
// 它接收一个 gin.Context 对象，该对象包含了关于请求和响应的所有信息。
// 在这个函数中，我们返回一个 JSON 响应，状态码为 200，消息为 "pong"。
func Ping(c *gin.Context) {
    // c.JSON 是一个方法，用于发送一个 JSON 响应。
    // 它接收两个参数：HTTP 状态码和要发送的数据。
    // gin.H 是一个快捷方式，用于创建一个 map[string]interface{}。
    c.JSON(200, gin.H{
        "message": "pong",
    })
}