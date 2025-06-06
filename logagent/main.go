package main

// logagent 是一个日志收集客户端，负责从指定位置收集的日志文件，并将日志发送到 Kafka
import (
	"fmt"
	"logagent/common"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"

	"gopkg.in/ini.v1"
)

//.ini 文件中的字段名通常是小写（如 address），而 Go 的结构体字段名通常是大写开头（如 Address），以符合 Go 的导出规则

// 整个程序的配置结构体
type Config struct {
	KafkaConfig `ini:"kafka"` // Kafka 配置
	EtcdConfig  `ini:"etcd"`  // Etcd 配置
}

type KafkaConfig struct {
	Address  string `ini:"address"`   // Kafka 地址 ini配置文件里该对应的字段为address
	Topic    string `ini:"topic"`     // Kafka 主题 ini配置文件里该对应的字段为topic
	ChanSize int    `ini:"chan_size"` // 消息通道大小 ini配置文件里该对应的字段为chansize
}
type EtcdConfig struct {
	Address    string `ini:"address"`     // Etcd 地址 ini配置文件里该对应的字段为address
	CollectKey string `ini:"collect_key"` // 日志收集配置项的键 ini配置文件里该对应的字段为collectkey
}

// etcd的写入可以通过 etcddemo/put get/main.go 中的代码来实现(注意根据ip修改写入的key)
// logagent日志收集客户端
// 功能：
// 1.根据部署logagent的服务器ip来监听etcd并从中收集日志文件配置项通过Tail追踪对应的日志文件，当日志内容发生增删改时发送到 Kafka 主题
// 2.当etcd中的日志收集配置项发生变化时，客户端通过Tail会自动更新配置并重新追踪新的日志文件
func main() {
	//-1:获取本机ip,为每个部署logagent的服务器后续从etcd对应服务器ip的键中的获取日志收集配置项做准备
	ip := common.GetOutboundIP() // 获取本机 IP 地址

	var configObj = new(Config) // 创建一个 Config 结构体实例
	// 0.读取配置文件 go ini
	// 使用 ini.MapTo 函数将配置文件内容映射到 Config 结构体
	// 注意：ini.MapTo 函数会根据结构体的标签（如 `ini:"kafka"`）来匹配配置文件中的字段
	err := ini.MapTo(configObj, "conf/config.ini")
	if err != nil {
		panic("Failed to load config file: " + err.Error())
	}
	fmt.Printf("%+v\n", configObj) // 打印配置对象
	// 1.连接到 Kafka
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize) // 传入 Kafka 地址
	if err != nil {
		panic("Failed to connect to Kafka: " + err.Error())
	}
	fmt.Println("Connected to Kafka successfully!")

	//2.初始化etcd连接
	err = etcd.Init([]string{configObj.EtcdConfig.Address}) // 传入 Etcd 地址
	if err != nil {
		panic("Failed to connect to Etcd: " + err.Error())

	}
	//从etcd中拉取需要收集的日志配置项（解决不能同时处理多个日志文件的问题）
	collectKey := fmt.Sprintf(configObj.EtcdConfig.CollectKey, ip) // 使用本机 IP地址替换collect_key 中的占位符(即格式化config.ini 中的 collect_key 键)
	allConf, err := etcd.GetConf(collectKey)                       // 传入日志收集配置项的键
	if err != nil {
		panic("Failed to get log collection config from Etcd: " + err.Error())

	}
	fmt.Println(allConf)
	//创建一个协程去监听 etcd 中的日志收集配置项变化
	go etcd.WatchConf(collectKey) // 传入日志收集配置项的键和初始配置

	// 3.tail 包来实现对每个日志文件的实时监控和读取并发送到 Kafka,
	err = tailfile.Init(allConf) // 传入日志文件路径
	if err != nil {
		panic("Failed to initialize tail file: " + err.Error())
	}
	fmt.Println("Tail file initialized successfully!")
	// 阻塞主协程
	select {}

}

// 以上代码实现了一个日志收集客户端，能够从指定位置收集日志文件，并将日志发送到 Kafka。
//使用了 etcd 来存储日志收集配置项，并使用 tailfile 包从etcd获取日志文件配置项并实时监控和读取对应日志文件内容发送到kafka。配置文件采用 ini 格式，包含 Kafka、Etcd
// common包里定义了日志收集配置项的结构体 CollectEntry 以及获取本机 IP 地址的函数 GetOutboundIP()。
