package main

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/shirou/gopsutil/cpu"
)

// influxdb go操作官方参考：https://docs.influxdb.org.cn/influxdb/v2/api-guide/client-libraries/go/
//以下是一个简单的influxdb demo程序 它每隔5秒获取一次CPU使用率，并数据写入到InfluxDB中

var client influxdb2.Client // 声明一个全局的 InfluxDB 客户端变量
// Get-CPU-Info 获取 CPU信息
func GetCpuInfo() float64 {
	// 获取 CPU 使用率
	cpuUsage, err := cpu.Percent(0, false) // 获取 CPU 使用率，false 表示不获取每个核心的使用率
	if err != nil {
		fmt.Println("获取 CPU 信息失败:", err) // 如果获取失败，打印错误信息
		return 0                         // 如果获取失败，返回错误
	}
	return cpuUsage[0] // 返回 CPU 使用率的第一个元素（表示总的 CPU 使用率）
}
func main() {

	// 创建 InfluxDB 客户端
	client = influxdb2.NewClient("http://192.168.219.129:8086", "X2U1L2pA73nMyiRbKaov9mFh247BlFG2iCukii2giKmqz0gJD4zUm4091mvORvCVgvxO39o-UN2x1z1bpu1drQ==")

	// 获取写入 API
	writeAPI := client.WriteAPIBlocking("logagent", "my-bucket")

	// 开始一个无限循环，每隔 5 秒获取一次 CPU 使用率并写入 InfluxDB
	for {
		time.Sleep(5 * time.Second) // 每隔 5 秒获取一次 CPU 使用率
		cpuUsage := GetCpuInfo()    // 调用 GetCpuInfo 函数获取 CPU 使用率
		// 创建一个点p 包含测量名称、标签、字段、时间戳
		measurement := "cpu_usage"
		tags := map[string]string{"cpu": "cpu0"}
		fields := map[string]interface{}{
			"cpu_percent": cpuUsage, // 将 CPU 使用率作为字段

		}

		p := influxdb2.NewPoint(measurement, tags, fields, time.Now())

		// 写入点到 InfluxDB
		if err := writeAPI.WritePoint(context.Background(), p); err != nil {
			fmt.Printf("Error writing point: %v\n", err)
		} else {
			fmt.Println("CPU usage written to InfluxDB:", cpuUsage)
		}
	}
}
