# scoreboard
使用Redis在服务器上对用户的点击数排序，并返回点击次数排行榜。



## 技术栈

Vue Gin Redis



## API设计

本项目API设计采用的是websocket实现。

由于考虑到用户在点击比较频繁，如果采用HTTP会造成头部开销较大，而websocket的头部开销会相对小一些。



## 前端

前端采用vue框架编写完成，UI组件采用elementplus


## 后端

```
backend
├── controllers
│   └── websocket.go
├── go.mod
├── go.sum
├── main.go
├── routers
│   └── router.go
└── services
    └── redis-server.go

```

后端采用Gin框架完成，大致流程：

- 在main.go中启动路由，并且启动端口监听
- 在routers/router.go中定义/ws路由，用于接收websocket的连接
- 对于ws的处理，函数定于在controllers/websocket.go中，包括针对不同任务类型使用redis数据库的函数调用
- 在services/redis-server.go中，对各个任务如何具体操作redis进行定义


## 运行方式

确保redis正常运行
``` shell
# 检查 Redis 是否正在运行
if ! pgrep -x "redis-server" > /dev/null
then
    # 如果 Redis 没有运行，就启动它
    redis-server &
fi
```

``` shell
cd backend
go run main.go
```

``` shell
cd frontend
npm install # 只需要在第一次运行时执行
npm run dev
```




