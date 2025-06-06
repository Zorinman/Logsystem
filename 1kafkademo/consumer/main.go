package main

import (
	"fmt"

	"sync"

	"github.com/Shopify/sarama"
)

// kafka 消费者示例
// ！！这里一定要使用waitgroup不然会导致goroutine 启动后还没来得及处理消息，main() 就退出了！
func main() {
	var wg sync.WaitGroup // 使用 WaitGroup防止main提前退出，等待所有 goroutine 完成
	// 1. 连接到 Kafka
	// 传入 Kafka 地址和配置使用 Sarama 库创建一个异步消费者（客户端）
	consumer, err := sarama.NewConsumer([]string{"192.168.219.132:9092"}, nil) //nil 表示使用默认配置
	if err != nil {
		fmt.Println("连接到 Kafka 出错:", err) // 如果连接失败，打印错误信息
		return
	}
	partitionList, err := consumer.Partitions("Test") // 获取主题 "Test" 的所有分区
	if err != nil {
		fmt.Println("获取分区列表出错:", err) // 如果获取分区列表失败，打印错误信息
		return
	}
	fmt.Println("分区列表:", partitionList) // 打印分区列表
	for partition := range partitionList {
		// 2. 订阅分区
		pc, err := consumer.ConsumePartition("Test", int32(partition), sarama.OffsetNewest) // 订阅主题 "Test" 的指定分区，使用最新的偏移量，也可以使用OffsetOldest 来从最早的消息开始消费
		if err != nil {
			fmt.Println("订阅分区出错:", err) // 如果订阅分区失败，打印错误信息
			return
		}
		fmt.Printf("开始消费分区 %d\n", partition) // 打印正在消费的分区信息

		// 3. 异步从每个分区消费消息
		wg.Add(1) // 增加 WaitGroup 的计数器，表示有一个新的 goroutine 开始工作
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			fmt.Println("开始处理分区的消息") // 打印正在处理的分区信息
			defer pc.AsyncClose()    // 在函数结束时关闭分区消费者
			for {
				select {
				case msg := <-pc.Messages(): // 从分区消费者中接收消息
					fmt.Printf("收到消息: 分区 %d, 偏移量 %d, 内容: %s\n", msg.Partition, msg.Offset, string(msg.Value)) // 打印消息内容
				case err := <-pc.Errors(): // 处理错误
					fmt.Println("消费错误:", err) // 如果有错误，打印错误信息
				}
			}
		}(pc) // 启动一个 goroutine 来处理每个分区的消息
	}
	wg.Wait()
}
