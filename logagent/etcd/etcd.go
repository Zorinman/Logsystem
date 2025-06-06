package etcd

import (
	"encoding/json"
	"fmt"
	"time"

	"logagent/common"
	"logagent/tailfile"

	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
)

var (
	client *clientv3.Client // 声明一个全局的 etcd 客户端变量
)

// Init 函数用于初始化 etcd 客户端
func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5, // 设置连接超时时间为 5 秒
	})
	if err != nil {
		panic(err) // 如果连接失败，抛出错误
	}
	fmt.Println("连接到 etcd 成功") // 打印连接成功的消息
	// defer client.Close()       // 在函数结束时关闭客户端连接
	return nil
}

// 拉取日志收集配置项
func GetConf(key string) (collectEntrylist []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // 设置上下文超时时间为 5 秒
	defer cancel()                                                          // 在函数结束时取消上下文，释放资源
	resp, err := client.Get(ctx, key)                                       // 从 etcd 中获取以 "collect" 为前缀的键值对
	if err != nil {
		fmt.Println("获取配置失败:", err) // 如果获取失败，打印错误信息
		return
	}
	if len(resp.Kvs) == 0 { // 检查是否有获取到的键值对
		fmt.Println("没有获取到配置") // 如果没有获取到配置，打印提示信息
		return
	}

	ret := resp.Kvs[0]
	fmt.Println(ret.Value)                             // 打印获取到的值
	err = json.Unmarshal(ret.Value, &collectEntrylist) // 将 JSON 格式的值解析到 collectEntrylist 列表中
	if err != nil {
		fmt.Println("解析配置失败:", err) // 如果解析失败，打印错误信息
		return
	}
	return

	//解答：cli.Get 获取"ip_collect_logs_conf"后 resp.Kvs[0] 就是键 "ip_collect_logs_conf" 的键值了，那么为什么还要用切片形式？
	//原因在于 Etcd 的设计允许通过一个请求返回多个键值对。虽然在您的场景中，cli.Get 只获取了一个键 "ip_collect_logs_conf"，但 Etcd 的 API 是通用的，支持多种查询模式
	// 	Etcd 的查询模式：
	// Etcd 支持多种查询模式，例如精确查询和前缀查询。
	// 精确查询：如果查询的是单个键（如 "ip_collect_logs_conf"），resp.Kvs 可能只包含一个元素。
	// 前缀查询：如果查询的是某个前缀（如 "ip_collect_logs_conf/"），resp.Kvs 可能包含多个键值对。
}

// WatchConf 函数用于监听 etcd 中的日志收集配置项变化
// key键中的内容发生变化（如添加、修改或删除）时，触发 WatchConf 函数
func WatchConf(key string) {
	for {

		Wchan := client.Watch(context.Background(), key)
		for wresp := range Wchan {
			for _, ev := range wresp.Events {
				var newConf []common.CollectEntry // 声明一个新的日志收集配置项列表
				fmt.Printf("收到新日志项: %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				//1.如果etcd中的键Key被删除（注意不是键的值被删除），则事件类型为删除事件，此时返回一个空的日志收集配置项列表
				//之后Tail会通知 tailfile 包更新日志收集配置为空
				if ev.Type == clientv3.EventTypeDelete {
					fmt.Printf("Etcd键%s被删除,日志配置项被清空！!:", string(ev.Kv.Key)) // 打印删除事件的键
					tailfile.SendNewConf(newConf)                           // 通知 tailfile 包更新日志收集配置为空
					continue
				}
				//2.如果etcd中键的值发生变化（如添加、删除、修改），则事件类型为 PUT，此时需要解析新的日志收集配置项
				err := json.Unmarshal(ev.Kv.Value, &newConf) // 将 JSON 格式的值解析到 newConf 中
				if err != nil {
					fmt.Println("解析新配置失败:", err) // 如果解析失败，打印错误信息
					continue
				}

				//通知 tailfile 包更新日志收集配置
				tailfile.SendNewConf(newConf)
			}
		}
	}

}
