package main

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// 将CPU信息写入到influxDB中
func writesCpuPoints(data *CpuInfo, cli influxdb2.Client) {
	// 获取写入 API
	writeAPI := cli.WriteAPIBlocking("logagent", "my-bucket")

	// 创建一个点p 包含测量名称、标签、字段、时间戳
	measurement := "cpu_usage"
	tags := map[string]string{"cpu": "cpu0"}
	fields := map[string]interface{}{
		"cpu_percent": data.CpuPercent,
	}

	p := influxdb2.NewPoint(measurement, tags, fields, time.Now())

	// 写入点p到 InfluxDB
	if err := writeAPI.WritePoint(context.Background(), p); err != nil {
		fmt.Printf("Error writing point: %v\n", err)
	} else {
		fmt.Println("CPU info cpu0 written to InfluxDB:")
	}
}

// 将内存信息写入到influxDB中
func writesMemPoints(data *MemInfo, cli influxdb2.Client) {

	// 获取写入 API
	writeAPI := cli.WriteAPIBlocking("logagent", "my-bucket")
	// 根据传入数据的类型插入数据
	measurement := "memory"
	tags := map[string]string{"mem": "mem"}
	fields := map[string]interface{}{
		"total":        int64(data.Total),
		"available":    int64(data.Available),
		"used":         int64(data.Used),
		"used_percent": data.UsedPercent,
	}

	p := influxdb2.NewPoint(measurement, tags, fields, time.Now())

	// 写入点p到 InfluxDB
	if err := writeAPI.WritePoint(context.Background(), p); err != nil {
		fmt.Printf("Error writing point: %v\n", err)
	} else {
		fmt.Println("Memory usage written to InfluxDB:")
	}
}

// 将磁盘信息写入到influxDB中
func writesDiskPoints(data *DiskInfo, cli influxdb2.Client) {
	// 获取写入 API
	writeAPI := cli.WriteAPIBlocking("logagent", "my-bucket")
	// 根据传入数据的类型插入数据
	measurement := "disk"
	for k, v := range data.PartitionUsageStat {
		tags := map[string]string{"path": k} // 使用分区的挂载点作为标签
		fields := map[string]interface{}{
			"total":               int64(v.Total),
			"free":                int64(v.Free),
			"used":                int64(v.Used),
			"used_percent":        v.UsedPercent,
			"inodes_total":        int64(v.InodesTotal),
			"inodes_used":         int64(v.InodesUsed),
			"inodes_free":         int64(v.InodesFree),
			"inodes_used_percent": v.InodesUsedPercent,
		}
		p := influxdb2.NewPoint(measurement, tags, fields, time.Now())
		// 写入点p到 InfluxDB
		if err := writeAPI.WritePoint(context.Background(), p); err != nil {
			fmt.Printf("Error writing point: %v\n", err)
		} else {
			fmt.Printf("disk info %v written to InfluxDB\n", k)
		}
	}

}

// 将网卡信息写入到influxDB中
func writesNetPoints(data *NetInfo, cli influxdb2.Client) {
	// 获取写入 API
	writeAPI := cli.WriteAPIBlocking("logagent", "my-bucket")
	// 根据传入数据的类型插入数据
	measurement := "net"
	for k, v := range data.NetIOCountersStat {
		tags := map[string]string{"name": k} // 把每个网卡存为tag
		fields := map[string]interface{}{
			"bytes_sent_rate":   v.BytesSentRate,
			"bytes_recv_rate":   v.BytesRecvRate,
			"packets_sent_rate": v.PacketsSentRate,
			"packets_recv_rate": v.PacketsRecvRate,
		}
		p := influxdb2.NewPoint(measurement, tags, fields, time.Now())
		// 写入点p到 InfluxDB
		if err := writeAPI.WritePoint(context.Background(), p); err != nil {
			fmt.Printf("Error writing point: %v\n", err)
		} else {
			fmt.Printf("Net info %v written to InfluxDB\n", k)
		}
	}
}
