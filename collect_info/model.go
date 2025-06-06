package main

import (
	"github.com/shirou/gopsutil/disk"
)

type CpuInfo struct {
	CpuPercent float64 `json:"cpu_percent"`
}

type MemInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Buffers     uint64  `json:"buffers"`
	Cached      uint64  `json:"cached"`
}

// type UsageStat struct {
// 	Path              string  `json:"path"`
// 	Fstype            string  `json:"fstype"`
// 	Total             uint64  `json:"total"`
// 	Free              uint64  `json:"free"`
// 	Used              uint64  `json:"used"`
// 	UsedPercent       float64 `json:"used_percent"`
// 	InodesTotal       uint64  `json:"inodes_total"`
// 	InodesUsed        uint64  `json:"inodes_used"`
// 	InodesFree        uint64  `json:"inodes_free"`
// 	InodesUsedPercent float64 `json:"inodes_used_percent"`
// }

// DiskInfo 包含了每个分区的使用情况，使用挂载点作为关键字存储
type DiskInfo struct {
	PartitionUsageStat map[string]*disk.UsageStat
}

type IOStat struct {
	BytesSent       uint64
	BytesRecv       uint64
	PacketsSent     uint64
	PacketsRecv     uint64
	BytesSentRate   float64 `json:"bytes_sent_rate"`   // 发送字节率
	BytesRecvRate   float64 `json:"bytes_recv_rate"`   // 接收字节率
	PacketsSentRate float64 `json:"packets_sent_rate"` // 发送数据包率
	PacketsRecvRate float64 `json:"packets_recv_rate"` // 接收数据包率
}

type NetInfo struct {
	NetIOCountersStat map[string]*IOStat
}
