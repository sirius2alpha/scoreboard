# scoreboard
使用Redis在服务器上对用户的点击数排序，并返回点击次数排行榜。
![image-20240121155707304](https://s2.loli.net/2024/01/21/cntLqdiyb9I3aer.png)

![image-20240121155736031](https://s2.loli.net/2024/01/21/N457iIYb82EQne6.png)


## 技术栈

![scoreboard技术栈](https://s2.loli.net/2024/01/21/cVJzxCyDFtOTjBL.png)

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



## 可以改进的地方

### 用户在登录的时候遇到相同用户名，会把他直接刷新



### 手机端自适应功能差，体验不好

- 手机在点击按钮的时候。会触发双击浏览器双击放大的功能，影响体验
- 手机端的网页有时候滑动不了
- 有时候手机端最上面的两个按钮会被浏览器的头部遮挡，但是又滑动不上去



### 对于只登录而没有点击的用户，排行榜中会保留下来，但不会清理

上一次的点击间隔和上次点击时间都不会刷新，后端是根据间隔时间清理用户，虽然可以保留，但是一直保存着也不是办法，可以设置一个单独的时长进行清理。
