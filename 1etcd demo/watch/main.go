package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://192.168.219.129:2379"},
		DialTimeout: time.Second * 5, // 设置连接超时时间为 5 秒
	})
	if err != nil {
		panic(err) // 如果连接失败，抛出错误
	}
	fmt.Println("连接到 etcd 成功") // 打印连接成功的消息
	defer cli.Close()          // 在函数结束时关闭客户端连接

	//watch 监听etcd的变化
	Wchan := cli.Watch(context.Background(), "zorin")
	for wresp := range Wchan {
		for _, ev := range wresp.Events {
			fmt.Printf("日志内容变更为: %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
