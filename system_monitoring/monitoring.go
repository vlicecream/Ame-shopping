package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
	"time"
)

var (
	cli client.Client // 在全局变量声明
	lastNetIOStatTimeStamp int64 // 上一次获取网络IO数据的时间点
	lastNetInfo *NetInfo // 上一次的网络数据
)

// 封装err报错信息
func errs(err error, str string) {
	if err != nil {
		log.Fatalf("%s err: %v", str, err)
	}
}

// 连接influxdb
func connInflux() (err error) {
	cli, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "admin",
		Password: "",
	})
	return
}

// 利用gopsutil包拿到CPU使用率
func getCpuInfo() {
	// var一个CpuInfo对象
	var cpuInfo = new(CpuInfo)
	percent, _ := cpu.Percent(time.Second, false)
	fmt.Printf("cpu percent:%v\n", percent)
	// 拿到传给wcp的参数
	cpuInfo.CpuPercent = percent[0]
	writesCpuPoints(cpuInfo)
}

// 利用gopsutil包拿到Mem数据
func getMemInfo() {
	// var一个CpuInfo对象
	var memInfo = new(MemInfo)
	vtMemoryStat, _ := mem.VirtualMemory()
	fmt.Printf("mem info:%v\n", memInfo)
	// 把属性一一对应
	memInfo.Used = vtMemoryStat.Used
	memInfo.Free = vtMemoryStat.Free
	memInfo.Total = vtMemoryStat.Total
	memInfo.Available = vtMemoryStat.Available
	// 将memInfo传入函数中
	writesMemPoints(memInfo)
}

// 利用gopsutil包拿到Disk数据
func getDiskInfo() {
	// var并初始化对象
	var diskInfo = &DiskInfo{
		PartitionUsageStat: make(map[string]*disk.UsageStat, 18),
	}
	parts, err := disk.Partitions(true)
	errs(err, "disk")
	// 循环硬盘分区
	for _, part := range parts {
		usageStat, err := disk.Usage(part.Mountpoint)
		if err != nil {
			fmt.Printf("get Partitions failed, err:%v\n", err)
			continue // 这个分区不行了，跳过下一个分区
		}
		diskInfo.PartitionUsageStat[part.Mountpoint] = usageStat
	}
	// 写入
	writesDiskPoints(diskInfo)
}

// 利用gopsutil包拿到Net数据，但是我们要直接数据没有用，所以我们要处理成网卡的速率
func getNetInfo() {
	var netInfo = &NetInfo{
		NetIOCountersStat: make(map[string]*IOStat, 8),
	}
	currentTimeStamp := time.Now().Unix()
	netIOs, err := net.IOCounters(true)

	if err != nil {
		fmt.Printf("get net io counters err:", err)
	}
	for _, netIO := range netIOs {
		var ioStat = new(IOStat)
		ioStat.BytesSent = netIO.BytesSent
		ioStat.BytesRecv = netIO.BytesRecv
		ioStat.PacketsSent = netIO.PacketsSent
		ioStat.PacketsRecv = netIO.PacketsRecv
		// 将具体网卡数据的ioStat变量添加到map中
		netInfo.NetIOCountersStat[netIO.Name] = ioStat
		// 开始计算网卡每秒速率
		if lastNetIOStatTimeStamp == 0 || lastNetInfo == nil {
			continue
		}
		intervel := currentTimeStamp - lastNetIOStatTimeStamp
		ioStat.BytesSentRate = (float64(ioStat.BytesSent) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].BytesSent))/float64(intervel)
		ioStat.BytesRecvRate = (float64(ioStat.BytesRecv) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].BytesRecv))/float64(intervel)
		ioStat.PacketsSentRate = (float64(ioStat.PacketsSent) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].PacketsSent))/float64(intervel)
		ioStat.PacketsRecvRate = (float64(ioStat.PacketsRecv) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].PacketsRecv))/float64(intervel)

	}
	// 更新全局记录的上一次采集网卡的时间点和网卡数据
	lastNetIOStatTimeStamp = currentTimeStamp
	lastNetInfo = netInfo
	// 发送到influxdb
	writesNetPoints(netInfo)
}


// run运行函数
func run(interval time.Duration) {
	ticker := time.Tick(interval)
	for _ = range ticker {
		getCpuInfo()
		getMemInfo()
		getDiskInfo()
		getNetInfo()
	}
}

func main() {
	err := connInflux()
	errs(err, "connInflux")
	run(time.Second)
}
