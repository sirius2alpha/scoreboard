#!/bin/bash

root_dir=$(dirname $(realpath $0))

if [ ! -f $root_dir/../config.toml ]; then
    echo "找不到配置文件 config.toml"
    exit 1
fi
source $root_dir/../config.toml


# 定义端口号
port=$backend_listen_port
if [ -z "$port" ]; then
    echo "端口号未定义。"
    exit 1
fi

# 检查端口8080是否被占用
pid=$(lsof -ti:$port)

# 判断端口是否被占用
if [ ! -z "$pid" ]; then
    echo "端口 $port 被以下进程占用:"
    ps -p $pid

    # 询问用户是否结束进程
    read -p "是否要结束这个进程? (y/n): " answer

    if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
        kill $pid
        echo "进程已被终止。"
    else
        echo "操作已取消。"
    fi
else
    echo "端口$port 未被占用。"
fi
