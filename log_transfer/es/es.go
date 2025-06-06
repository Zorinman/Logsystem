package es

import (
	"bytes"
	"encoding/json"
	"fmt"

	elastic "github.com/elastic/go-elasticsearch/v7"
)

type ESClient struct {
	client      *elastic.Client  // Elasticsearch 客户端
	logDataChan chan interface{} // 通道，用于接收 Kafka 消费的日志数据
	index       string           // Elasticsearch 索引名称
}

var (
	esClient = new(ESClient) // 全局变量，存储 Elasticsearch 客户端实例
)

// Elasticsearch的连接和数据写入操作

func Init(addr []string, goroutineNum int, index string, maxSize int) (err error) {
	// Create a new Elasticsearch client
	esClient.client, err = elastic.NewClient(elastic.Config{ //初始化一个新的Elasticsearch客户端
		Addresses: addr,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Elasticsearch client created successfully")

	esClient.logDataChan = make(chan interface{}, maxSize) // 初始化通道
	esClient.index = index                                 // 设置索引名称
	for i := 0; i < goroutineNum; i++ {
		go sendToES(index) // 启动多个 goroutine 来处理日志数据的发送

	}
	return nil
}

// senToES 函数用于将通道中的日志数据写入到 Elasticsearch，没有消息则阻塞等待
// 用go关键字启动该函数的协程会持续阻塞监听 esClient.logDataChan 通道中的消息
func sendToES(index string) {

	for msg := range esClient.logDataChan {
		b, err := json.Marshal(msg) // 将日志数据序列化为 JSON 格式的字节切片
		if err != nil {
			fmt.Println("Error marshaling log data:", err)
			continue // 如果序列化失败，跳过当前循环
		}
		_, err = esClient.client.Index(
			index,

			bytes.NewReader(b), // 直接使用 bytes.NewReader(data) 直接将JSON 字节切片转换为 io.Reader
		)
		if err != nil {
			panic(err)
		}
		fmt.Println("Document indexed successfully")
	}

}

// 将kafka消费的日志数据发到通道中,以便后续处理
func PutLogData(m1 interface{}) error {
	esClient.logDataChan <- m1 // 将日志数据发送到通道中
	return nil                 // 返回 nil 表示操作成功
}
