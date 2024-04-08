#!/bin/bash

root_dir=$(dirname $(realpath $0))

if [ ! -f $root_dir/../conf/service.ini ]; then
    echo "找不到配置文件 service.ini"
    exit 1
fi
source $root_dir/../conf/service.ini

# 定义端口号
backend_listen_port=$backend_listen_port
vite_port=$vite_port

# port遍历
for port in $backend_listen_port $vite_port; do
    # 判断端口号是否定义
    if [ -z "$port" ]; then
        echo "端口号未定义。"
        exit 1
    fi

    # 检查端口是否被占用
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

done
