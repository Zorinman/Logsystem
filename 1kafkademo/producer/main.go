package main

//kafka 示例代码展示了如何使用 Sarama 库连接到 Kafka，并发送一条消息到指定的主题
import (
	"fmt"

	"github.com/Shopify/sarama"
)

func main() {
	//1. 创建生产者配置示例
	config := sarama.NewConfig()                              // 使用 Sarama 库创建一个新的配置对象
	config.Producer.RequiredAcks = sarama.WaitForAll          // 等待所有同步副本确认消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 使用随机分区器
	config.Producer.Return.Successes = true                   // 等待消息成功发送的确认

	//2. 连接到 Kafka
	// 传入 Kafka 地址和配置使用 Sarama 库创建一个同步生产者（客户端），同步生产者：在发送消息时会等待 Kafka 的确认，确保消息已成功写入 Kafka 的分区，这里client表示同步生产者
	client, err := sarama.NewSyncProducer([]string{"192.168.219.132:9092"}, config)
	if err != nil {
		fmt.Println("连接到 Kafka 出错:", err)
		return
	}
	defer client.Close() // 使用完生产者后关闭

	//3. 创建消息示例
	msg := &sarama.ProducerMessage{}
	msg.Topic = "Test"                                     // 主题名称
	msg.Value = sarama.StringEncoder("this is a test log") // 消息内容

	//4. 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("发送消息出错:", err)
		return
	} else {
		fmt.Printf("消息发送成功! 分区 ID: %d, 偏移量: %d\n", pid, offset)
	}

}
