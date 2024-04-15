# scoreboard
项目部署在：http://redisboard.sirius1y.top/

对于项目更详细的介绍，移步：[sirius1y.me](https://sirius1y.me/posts/notes/dev/dev-%E7%82%B9%E5%87%BB%E6%8E%92%E8%A1%8C%E6%A6%9C/)

使用Redis在服务器上对用户的点击数排序，并返回点击次数排行榜。
![网站首页](https://s2.loli.net/2024/01/21/cntLqdiyb9I3aer.png)

![排行榜](https://s2.loli.net/2024/01/21/CUjLTrbN9KPZ4xo.png)


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
make build # 构建项目
make run # 运行项目
```