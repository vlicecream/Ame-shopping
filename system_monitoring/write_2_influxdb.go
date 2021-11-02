package main

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"time"
)

// 把error封装成函数，避免冗余代码，影响可读性
func errors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


// 拿到cpu数据发送给influxdb数据库
func writesCpuPoints(data *CpuInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "ame", // 选择库
		Precision: "s",   //精度，默认ns
	})
	errors(err)
	tags := map[string]string{"cpu": "cpu"} // 设置标签
	fields := map[string]interface{}{ // 采集的信息
		"cpu_percent": data.CpuPercent,
	}

	pt, err := client.NewPoint("cpu_percent", tags, fields, time.Now())
	errors(err)
	bp.AddPoint(pt)

	err = cli.Write(bp)
	errors(err)
	log.Println("CPU insert success")
}

// 拿到Mem数据发送给influxdb数据库
func writesMemPoints(data *MemInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "ame", // 选择库
		Precision: "s",   //精度，默认ns
	})
	errors(err)
	tags := map[string]string{"Mem": "mem"} // 设置标签
	fields := map[string]interface{}{ // 采集的信息
		"total":     int64(data.Total),
		"available": int64(data.Available),
		"used":      int64(data.Used),
		"free":      int64(data.Free),
	}

	pt, err := client.NewPoint("Memory", tags, fields, time.Now())
	errors(err)
	bp.AddPoint(pt)

	err = cli.Write(bp)
	errors(err)
	log.Println("Mem insert success")
}

// 拿到Disk数据发送给influxdb数据库
func writesDiskPoints(data *DiskInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "ame", // 选择库
		Precision: "s",   //精度，默认ns
	})
	errors(err)

	// 循环遍历去除map中的key:value
	for k, v := range data.PartitionUsageStat {
		tags := map[string]string{"path": k} // 设置标签
		fields := map[string]interface{}{ // 采集的信息
			"total": int64(v.Total),
			"free":  int64(v.Free),
			"user":  int64(v.Used),
		}
		pt, err := client.NewPoint("Disk", tags, fields, time.Now())
		errors(err)
		bp.AddPoint(pt)
	}
	err = cli.Write(bp)
	errors(err)
	log.Println("Disk insert success")
}

// 拿到Net数据发送给influxdb数据库
func writesNetPoints(data *NetInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "ame",
		Precision: "s", //精度，默认ns
	})
	errors(err)
	// 循环在map取出key:value 跟disk思路一致
	for k, v := range data.NetIOCountersStat {
		tags := map[string]string{"name": k} // 把每个网卡按名字建立索引
		fields := map[string]interface{}{
			"bytesSentRate": v.BytesSentRate,
			"bytesRecvRate": v.BytesRecvRate,
			"packetsSentRate": v.PacketsSentRate,
			"packetsRecvRate": v.PacketsRecvRate,
		}
		pt, err := client.NewPoint("net", tags, fields, time.Now())
		errors(err)
		bp.AddPoint(pt)
	}
	err = cli.Write(bp)
	errors(err)
	log.Println("Net insert success")
}