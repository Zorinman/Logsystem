// colletc_info 实现了从系统中收集CPU、内存、磁盘和网络IO信息，并将这些信息写入到InfluxDB数据库中
package main

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

var (
	cli                    influxdb2.Client
	lastNetIOStatTimeStamp int64    // 上一次获取网络IO数据的时间点
	lastNetInfo            *NetInfo // 上一次的网络IO数据
)

func run(interval time.Duration) {
	ticker := time.Tick(interval)
	for _ = range ticker {
		getCpuInfo()
		getMemInfo()
		getDiskInfo()
		getNetInfo()
	}
}

func initConnInflux() {
	// 创建 InfluxDB 客户端
	cli = influxdb2.NewClient("http://192.168.219.129:8086", "X2U1L2pA73nMyiRbKaov9mFh247BlFG2iCukii2giKmqz0gJD4zUm4091mvORvCVgvxO39o-UN2x1z1bpu1drQ==")

}
func getCpuInfo() {
	var cpuInfo = new(CpuInfo)
	// CPU使用率
	percent, _ := cpu.Percent(time.Second, false)
	// 写入到influxDB中
	cpuInfo.CpuPercent = percent[0]
	writesCpuPoints(cpuInfo, cli)
}
func getMemInfo() {
	var memInfo = new(MemInfo)
	info, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("get mem info failed, err:%v", err)
		return
	}
	memInfo.Total = info.Total
	memInfo.Available = info.Available
	memInfo.Used = info.Used
	memInfo.UsedPercent = info.UsedPercent
	memInfo.Buffers = info.Buffers
	memInfo.Cached = info.Cached
	writesMemPoints(memInfo, cli)
}

// 遍历每个分区，通过分区的挂载点获取分区的使用情况，将每个分区的使用情况按照挂载点关键字存储到diskInfo中
func getDiskInfo() {
	var diskInfo = &DiskInfo{
		PartitionUsageStat: make(map[string]*disk.UsageStat, 16),
	}
	parts, _ := disk.Partitions(true)
	for _, part := range parts {
		// 拿到每一个分区
		usageStat, err := disk.Usage(part.Mountpoint) // 传挂载点
		if err != nil {
			fmt.Printf("get %s usage stat failed\n", part.Mountpoint)
			continue
		}
		diskInfo.PartitionUsageStat[part.Mountpoint] = usageStat
	}
	writesDiskPoints(diskInfo, cli)
}

// 获取网卡的字节发生和接收值以及包的发生和接收值，计算速率存储到netInfo，将netInfo中的速率写入到influxDB中
// 第一次调用getNetInfo()，for循环只获取网卡的IO数据，从第二次调用开始，for循环计算速率
func getNetInfo() {
	// 获取网卡IO数据,8个网卡,网卡名作为key
	var netInfo = &NetInfo{
		NetIOCountersStat: make(map[string]*IOStat, 8),
	}
	currentTimeStamp := time.Now().Unix()
	netIOs, err := net.IOCounters(true)
	if err != nil {
		fmt.Printf("get net io counters failed, err:%v", err)
		return
	}
	for _, netIO := range netIOs {
		var ioStat = new(IOStat)
		ioStat.BytesSent = netIO.BytesSent
		ioStat.BytesRecv = netIO.BytesRecv
		ioStat.PacketsSent = netIO.PacketsSent
		ioStat.PacketsRecv = netIO.PacketsRecv
		// 将具体网卡数据的ioStat变量添加到map中
		netInfo.NetIOCountersStat[netIO.Name] = ioStat // 不要放到continue下面

		// 开始计算网卡相关速率
		// 如果上一次采集网卡的时间点为0或者上一次的网卡数据为nil,则跳过计算速率
		if lastNetIOStatTimeStamp == 0 || lastNetInfo == nil {
			continue
		}
		// 计算时间间隔：当前时间戳 - 上一次采集网卡的时间戳
		interval := currentTimeStamp - lastNetIOStatTimeStamp
		// 计算速率：速率 = (当前采集的值 - 上一次采集的值) / 时间间隔
		ioStat.BytesSentRate = (float64(ioStat.BytesSent) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].BytesSent)) / float64(interval)
		ioStat.BytesRecvRate = (float64(ioStat.BytesRecv) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].BytesRecv)) / float64(interval)
		ioStat.PacketsSentRate = (float64(ioStat.PacketsSent) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].PacketsSent)) / float64(interval)
		ioStat.PacketsRecvRate = (float64(ioStat.PacketsRecv) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].PacketsRecv)) / float64(interval)

	}
	// 更新全局记录的上一次采集网卡的时间点和网卡数据
	lastNetIOStatTimeStamp = currentTimeStamp // 更新时间
	lastNetInfo = netInfo
	// 插入数据到influxDB
	writesNetPoints(netInfo, cli)
}

func main() {

	initConnInflux()

	run(time.Second * 5)
}
