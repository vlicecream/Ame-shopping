package main

import "github.com/shirou/gopsutil/disk"

// 定义系统信息类型常量
const (
	CpuInfoType  = "cpu"
	MemInfoType  = "mem"
	DiskInfoType = "disk"
	NetInfoType  = "net"
)

// SysInfo 定义一个概括系统信息的结构体
type SysInfo struct {
	InfoType string
	IP       string
	Data     interface{} // 因为类型不确定 用空接口
}

// CpuInfo 定义一个Cpu结构体
type CpuInfo struct {
	CpuPercent float64 `json:"cpu_percent"`
}

// MemInfo 定义一个Mem结构体
type MemInfo struct {
	// Total 系统中RAM的总量 RAM: 随机存取存储器 也叫主存
	Total uint64 `json:"total"`
	// 可用于程序分配的RAM
	Available uint64 `json:"available"`
	// 程序使用的RAM
	Used uint64 `json:"used"`
	// 内核的空余内存
	Free uint64 `json:"free"`
}

// DiskInfo 定义一个Disk结构体
type DiskInfo struct {
	PartitionUsageStat map[string]*disk.UsageStat
}

// IOStat 定义一个IOStat类型
type IOStat struct {
	BytesSent   uint64 `json:"bytesSent"`   // number of bytes sent
	BytesRecv   uint64 `json:"bytesRecv"`   // number of bytes received
	PacketsSent uint64 `json:"packetsSent"` // number of packets sent
	PacketsRecv uint64 `json:"packetsRecv"` // number of packets received
	BytesSentRate   float64 `json:"bytesSentRate"`   // number of bytes sent
	BytesRecvRate  float64 `json:"bytesRecvRate"`   // number of bytes received
	PacketsSentRate float64 `json:"packetsSentRate"` // number of packets sent
	PacketsRecvRate float64 `json:"packetsRecvRate"` // number of packets received
}

// NetInfo 把IOStat封装成map 方便管理
type NetInfo struct {
	NetIOCountersStat map[string]*IOStat
}