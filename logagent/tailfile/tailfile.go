package tailfile

import (
	"encoding/json"
	"fmt"
	"logagent/common"
	"logagent/kafka"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
)

type tailTask struct {
	path  string     // 日志文件路径
	topic string     // 日志主题
	tObj  *tail.Tail // tail.Tail 的实例
}

// LogMessage 结构体用于存储最终发送到 Kafka 的日志消息
type LogMessage struct {
	IP       string `json:"ip"`
	Path     string `json:"path"`
	LineText string `json:"linetext"` // 日志文件每行内容
}

// Run方法用于启动 tailTask 实例，开始追踪日志文件
func (tt *tailTask) Run() {
	// 该方法用于启动 tailTask 实例，开始追踪日志文件
	// 从 tt.tObj 中读取新行，发送到 Kafka
	for line := range tt.tObj.Lines { // 从 tt.tObj 中读取新行
		if line.Err != nil { // 检查是否有错误
			fmt.Println("Error reading line:", line.Err)
			time.Sleep(time.Second * 1) // 如果有错误，等待 1 秒后重试
			continue                    // 如果有错误，跳过当前循环
		}

		if len(strings.Trim(line.Text, "\r")) == 0 { // 检查行是否为空,使用strings.Trim函数去除行首尾的空白回车符
			continue // 如果行为空，跳过当前循环
		}
		// 对每一行日志内容进行处理
		// 创建一个 LogMessage 实例进行封装，包含 IP、Path 和 LineText
		logMsg := LogMessage{
			IP:       common.GetOutboundIP(), // 获取本机IP地址
			Path:     tt.path,
			LineText: line.Text,
		}
		// 将日志消息转换为 JSON 格式
		jsonBytes, err := json.Marshal(logMsg)
		if err != nil {
			fmt.Println("JSON 编码失败:", err)
			continue
		}

		// msg := &sarama.ProducerMessage{ // 创建一个新的消息
		// 	Topic: tt.topic,                        // 设置主题
		// 	Value: sarama.StringEncoder(line.Text), // 设置消息内容
		// }
		//利用MessageChannel异步地将日志发送到Kafka
		msg := &sarama.ProducerMessage{ // 创建一个新的消息
			Topic: tt.topic,                      // 设置主题
			Value: sarama.ByteEncoder(jsonBytes), // 设置消息内容
		}
		kafka.SendMsgChan(msg)                                                                           // 将消息发送到 Msgchan 通道
		fmt.Printf("新日志内容已发送到msgchan等待发送到kfaka....:日志路径:%s,主题:%s,内容:%s\n", tt.path, tt.topic, line.Text) // 打印日志收集成功的消息
	}

	// Meschan在kafka/kafka.go中被SendMsg()函数读取，并发送到Kafka
}

// newTailTask 函数用于创建一个新的 tailTask 实例
func newTailTask(Path, Topic string) (*tailTask, error) {
	var err error // 声明一个错误变量，用于存储可能发生的错误
	cfg := tail.Config{
		ReOpen:    true,                                 // 如果文件被轮转（如日志切割），重新打开文件
		Follow:    true,                                 // 跟随文件的增长（类似于 tail -f）
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件末尾开始读取
		Poll:      true,                                 // 使用轮询模式检测文件变化
		MustExist: true,                                 // 确保文件存在
	}
	tt := tailTask{
		path:  Path,  // 日志文件路径
		topic: Topic, // 日志主题
	}
	tt.tObj, err = tail.TailFile(Path, cfg) // 使用指定的cfg配置创建 tail.Tail 实例
	if err != nil {
		return nil, err // 如果创建失败，返回 nil
	}
	return &tt, nil
	// 返回一个指向 tailTask 的指针
}

// func Init(allConf []common.CollectEntry) (err error) {

// 	//allConf存了若干日志收集项
// 	//针对每一个日志收集项，创建一个 tailTask 实例
// 	for _, conf := range allConf {
// 		tt, err := newTailTask(conf.Path, conf.Topic) // 创建一个新的 tailTask 实例
// 		if err != nil {
// 			fmt.Printf("路径%s文件tailTask追踪实例创建失败:%s", conf.Path, err)
// 			return err // 如果创建失败，返回错误
// 		}
// 		fmt.Println("为conf.Path创建tailTask追踪实例成功")    // 打印成功消息
// 		go tt.Run()                                  // 启动 tailTask 实例，开始追踪日志文件
// 		fmt.Printf("路径%s文件Run监听线程创建成功\n", conf.Path) // 打印成功消息

// 	}
// 	//初始化用来接收配置项的通道
// 	confChan = make(chan []common.CollectEntry) // 创建阻塞通道 confChan，用于接收新的日志收集配置项
// 	//等待etcd的WatchConf通知有新的日志收集项

// 	newConf := <-confChan                            // 从 confChan 通道中接收新的日志收集配置项
// 	fmt.Printf("Tail收到新的日志收集配置项%v，开始处理...", newConf) // 打印收到新配置的消息

// 	//新配置到达之后

// 	return
// }
