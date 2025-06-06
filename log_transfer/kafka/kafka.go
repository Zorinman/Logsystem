package kafka

import (
	"encoding/json"
	"fmt"
	"logtransfer/es"

	"github.com/Shopify/sarama"
)

// 初始化kafka连接
// 消费kafka的日志消息，写入到elasticsearch中

func Init(addr []string, topic string) (err error) {
	// 1. 连接到 Kafka
	// 传入 Kafka 地址和配置使用 Sarama 库创建一个异步消费者（客户端）
	consumer, err := sarama.NewConsumer(addr, nil) //nil 表示使用默认配置
	if err != nil {
		fmt.Println("连接到 Kafka 出错:", err) // 如果连接失败，打印错误信息
		return
	}
	partitionList, err := consumer.Partitions(topic) // 获取主题 的所有分区
	if err != nil {
		fmt.Println("获取分区列表出错:", err) // 如果获取分区列表失败，打印错误信息
		return
	}
	for partition := range partitionList {
		// 2. 订阅分区
		var pc sarama.PartitionConsumer                                                   // 声明一个分区消费者变量，用于异步消费消息
		pc, err = consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest) // 订阅主题 addr 的指定分区，使用最新的偏移量，也可以使用OffsetOldest 来从最早的消息开始消费
		if err != nil {
			fmt.Println("订阅分区出错:", err) // 如果订阅分区失败，打印错误信息
			return err
		}
		fmt.Printf("开始消费分区 %d\n", partition) // 打印正在消费的分区信息

		// 3. 异步从每个分区消费消息
		//pc.Messages() 返回的是通道，因此协程会阻塞在这里，持续监听kafka中的新消息
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Println(msg.Topic, string(msg.Value)) // 打印消息内容
				var m1 map[string]interface{}             // 声明一个 map 用于存储消息内容
				err = json.Unmarshal(msg.Value, &m1)      // 将消息内容解析并存储到 m1中
				if err != nil {
					fmt.Println("解析消息内容出错:", err)
					continue // 如果解析失败，跳过当前循环
				}
				fmt.Println("解析后的消息内容:", m1)
				es.PutLogData(m1) // 将解析后的消息内容写入到 Elasticsearch 中

			}

		}(pc) // 启动一个 goroutine 来处理每个分区的消息
	}
	return nil // 返回 nil 表示初始化成功
}
