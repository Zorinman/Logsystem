package main

import (
	"fmt"
	"logtransfer/es"
	"logtransfer/kafka"
	"logtransfer/model"

	"gopkg.in/ini.v1"
)

// 从kafka中消费日志数据，写入到elasticsearch中
// kafka.go中 msg.values是从kafka中取出的原始消息为[]byte类型,最后经过一系列处理写入到es中的bytes.NewReader(b)中b也是[]byte类型
// 处理过程为:`msg.value` -> `json.Unmarshal(msg.Value, &m1)` -> `es.PutLogData(m1)` → `esClient.logDataChan <- m1 `
// -> `for msg := range esClient.logDataChan {b, err := json.Marshal(msg)}` → `esClient.client.Index(index, bytes.NewReader(b))`
// 之所以反序列化Unmarshal又序列化Marshal有以下原因:
// 1. ⭐验证 Kafka 消息是否为合法 JSON格式：Kafka 消息是原始[]byte，[]byte中不一定是标准的 JSON 格式，非JSON格式经过json.Unmarshal会报错，过滤了非JSON格式的日志数据
// 2. 支持对日志结构进行加工或增强

func main() {
	var cfg = new(model.Config)
	err := ini.MapTo(cfg, "config/logtransfer.ini")
	if err != nil {
		panic("Failed to load config file: " + err.Error())
	}
	fmt.Println("Config loaded successfully:", cfg)
	//顺序不能颠倒，先启动Elasticsearch，再启动到Kafka，使得es.PutLogData()中的通道可以正常工作
	//1.连接到Elasticsearch
	err = es.Init([]string{cfg.EsConf.Address}, cfg.EsConf.GoroutineNum, cfg.EsConf.Index, cfg.EsConf.ChanSize) // 传入 Elasticsearch 地址、协程数、索引名称和通道大小
	if err != nil {
		panic("Failed to connect to Elasticsearch: " + err.Error())
	}
	fmt.Println("Connected to Elasticsearch successfully!")
	// 2.连接到 Kafka
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.Topic) // 传入 Kafka 地址和主题
	if err != nil {
		panic("Failed to connect to Kafka: " + err.Error())
	}
	fmt.Println("Connected to Kafka successfully!")

	select {} // 阻塞主 goroutine，防止程序退出
}
