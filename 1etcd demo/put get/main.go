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

	//put 将值放入etcd
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)                          // 设置上下文超时时间为 5 秒
	str := `[{"path":"log/test.log","topic":"web_logs"},{"path":"log/test2.log","topic":"web_logs"}]` // 要存储的 JSON 字符串
	_, err = cli.Put(ctx, "192.168.0.101_collect_logs_conf", str)
	if err != nil {
		fmt.Println("写入键值对失败:", err) // 如果写入失败，打印错误信息
		return

	}
	fmt.Println("写入键值对成功") // 如果写入成功，打印成功信息
	cancel()               // 取消上下文，释放资源

	//get 从etcd取值
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10) // 设置上下文超时时间为 5 秒
	resp, err := cli.Get(ctx, "zorin")                                      // 从 etcd 中获取键 "zorin" 的值
	if err != nil {
		fmt.Println("获取键值对失败:", err) // 如果获取失败，打印错误信息
		return
	}
	cancel()                      // 取消上下文，释放资源
	for _, kv := range resp.Kvs { // 遍历获取到的键值对
		fmt.Printf("键: %s, 值: %s\n", kv.Key, kv.Value) // 打印键和值
	}

}
