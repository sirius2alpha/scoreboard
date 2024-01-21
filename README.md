# scoreboard

对于项目更详细的介绍，移步：https://sirius1y.me/posts/notes/dev/dev-%E7%82%B9%E5%87%BB%E6%8E%92%E8%A1%8C%E6%A6%9C/



使用Redis在服务器上对用户的点击数排序，并返回点击次数排行榜。
![image-20240121155707304](https://s2.loli.net/2024/01/21/cntLqdiyb9I3aer.png)

![排行榜](https://s2.loli.net/2024/01/21/CUjLTrbN9KPZ4xo.png)


## 技术栈

![scoreboard技术栈](https://s2.loli.net/2024/01/21/cVJzxCyDFtOTjBL.png)

## 整体设计

1. **用户界面**
   排行榜展示区: 显示当前排行榜的状态。
   点击按钮: 用户点击来增加他们的计分。
   昵称输入和提交: 允许新用户输入昵称并参与排行榜。
   实时更新监听: 不需要用户交互，自动更新排行榜。
2. **WebSocket客户端逻辑**
   建立连接: 当用户访问网站时，建立WebSocket连接。
   发送点击事件: 当用户点击按钮时，发送消息到服务器。
   接收排行榜更新: 监听来自服务器的排行榜更新，并更新界面。
   用户注册: 发送新用户的昵称到服务器。
   处理断开连接: 如果用户20秒未操作，发送断开消息到服务器。
   后端设计（Gin + Redis）
3. **WebSocket服务器**
   处理WebSocket连接: 接受和管理WebSocket连接。
   接收消息: 解析从客户端接收到的消息（点击事件，新用户注册）。
   Redis交互: 更新用户的分数并重新排序排行榜。
   广播排行榜更新: 将更新后的排行榜发送给所有连接的客户端。
   处理断开: 移除30秒未操作的用户。
4. **Redis逻辑**
   用户分数管理: 存储和更新用户分数。
   排行榜排序: 实时更新排行榜。
   数据持久化: 保证数据在服务重启后仍然可用。



## API设计

本项目API设计采用的是websocket实现。

由于考虑到用户在点击比较频繁，如果采用HTTP会造成头部开销较大，而websocket的头部开销会相对小一些。

#### 消息类型

- UserClick: { type: "UserClick", nickname: "用户昵称" }

- NewUser: { type: "NewUser", nickname: "用户昵称" }

- UserInactive: { type: "UserInactive", nickname: "用户昵称" }

- RankUpdate: { type: "RankUpdate", ranks: [{nickname: "用户昵称", score: 分数,ClickTime: 上次点击时间, ClickInterval: 上次点击间隔时间}, ...] }

#### API流程

用户点击: 前端发送UserClick消息到服务器。
新用户加入: 前端发送NewUser消息到服务器。
服务器处理: 接收消息，更新Redis数据，并重新排序排行榜。
排行榜更新: 服务器广播RankUpdate消息到所有客户端。
前端更新界面: 客户端接收RankUpdate消息，更新排行榜显示。



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
