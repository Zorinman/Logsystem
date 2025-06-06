package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var msgchan chan *sarama.ProducerMessage // 声明一个全局变量 msgchan 用于存储消息通道 不对外暴露
var client sarama.SyncProducer           // 声明一个全局变量 client 用于存储 Kafka 的同步生产者
// kafka相关操作
// 和init不同,Init只是一个普通的自定义函数
func Init(address []string, chansize int) (err error) {
	//kafka初始化
	//address: kafka地址

	//1. Producer example
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // Wait for all in-sync replicas to acknowledge the message
	config.Producer.Partitioner = sarama.NewRandomPartitioner //partitioner
	config.Producer.Return.Successes = true                   // Wait for successful message delivery

	//2. connect to kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		fmt.Println("Error connecting to Kafka:", err)
		return err
	}
	msgchan = make(chan *sarama.ProducerMessage, chansize) // 创建一个消息通道，缓冲区大小为1000
	go SendMsg()                                           // 启动一个 goroutine 来处理消息发送
	return

}

// 从Mesgchan中读取消息并发送到Kafka
func SendMsg() {

	for {
		select {
		case msg := <-msgchan: // 从消息通道中读取消息
			pid, offset, err := client.SendMessage(msg) // 发送消息到 Kafka
			if err != nil {
				fmt.Println("Error sending message to Kafka:", err)
				continue // 如果发送失败，跳过当前循环
			}
			fmt.Printf("新日志内容成功发送到kafka! Partition ID: %d, Offset: %d\n", pid, offset) // 打印成功信息
		}

	}
}

// 这里mgschan用了小写不对外暴露，防止其它包直接访问 msgchan 通道并从中读取消息
// 外部仅可以调用SendMsgChan 对 msgchan 通道发消息
// 确保只有kafka包内 SendMsg 函数可以从 msgchan 中读取消息
func SendMsgChan(msg *sarama.ProducerMessage) {
	// 将消息发送到 msgchan 通道
	msgchan <- msg

}
