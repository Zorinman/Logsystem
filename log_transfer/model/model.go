package model

type Config struct {
	KafkaConf `ini:"kafka"`
	EsConf    `ini:"es"`
}
type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}
type EsConf struct {
	Address      string `ini:"address"`
	Index        string `ini:"index"`
	ChanSize     int    `ini:"chan_size"`     // 通道大小
	GoroutineNum int    `ini:"goroutine_num"` // 启动的协程数量
}
