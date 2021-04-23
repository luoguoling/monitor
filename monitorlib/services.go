//对服务层面的应用进行监控
package monitorlib

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"mwmonitor/config"
	newlog "mwmonitor/logger"
	"strconv"
	"time"
)

//float转化为string
func FloatToString(input float64) string {
	return strconv.FormatFloat(float64(input), 'f', 0, 64)
}

//监控cpu使用率
func MonitorCpu() {
	for {
		per, _ := cpu.Percent(time.Duration(config.GetConfig().Interval)*time.Second, false)
		if per[0] > config.GetConfig().AlterLimit {
			newlog.Mylog("系统报警").Warn("cpu使用率过高!!!!")
			SendDingMsg("项目cpu监控报警")
		}
	}

}

//监控cpu负载
func checkLoad() {
	cpuInfos, _ := cpu.Info()
	cpuNum := len(cpuInfos)
	loadinfo, _ := load.Avg()
	load5Info := loadinfo.Load5
	if int(load5Info) > 2*cpuNum {
		newlog.Mylog("系统报警").Warn("cpu负载报警!!!!")
		msg := "cpu负载过高已经超过警告值" + " " + FloatToString(config.GetConfig().AlterLimit) + "当前值:" + FloatToString(load5Info)
		SendDingMsg(msg)
	}
}
func MonitorLoad() {
	for {
		checkLoad()
		time.Sleep(time.Duration(config.GetConfig().Interval) * time.Second)
		//msg := "负载过高已经超过警告值" + " " + FloatToString(config.GetConfig().AlterLimit)
		//newlog.Mylog("系统报警").Warn(msg)
		//SendDingMsg(msg)
	}

}

//监控内存
func MonitorMem() {
	for {
		m, _ := mem.VirtualMemory()
		if m.UsedPercent > config.GetConfig().AlterLimit {
			msg := "内存报警:已经超过警告值" + " " + FloatToString(config.GetConfig().AlterLimit) + "当前值:" + FloatToString(m.UsedPercent)
			newlog.Mylog("系统报警").Warn(msg)
			SendDingMsg(msg)
		}
		time.Sleep(time.Duration(config.GetConfig().Interval) * time.Second)
	}
}

//监控磁盘
func checkDisk() []map[string]float64 {
	diskInfoList := make([]map[string]float64, 0)
	warnDiskList := make([]map[string]float64, 0)
	parts, _ := disk.Partitions(true)
	for _, part := range parts {
		pm := part.Mountpoint
		usage, _ := disk.Usage(pm)
		usagePercent := usage.UsedPercent
		//fmt.Println(usagePercent)
		diskInfoMap := make(map[string]float64)
		diskInfoMap[pm] = usagePercent
		//diskInfoList := make([]map[string]int,0)
		diskInfoList = append(diskInfoList, diskInfoMap)
	}
	for _, j := range diskInfoList {
		for k, v := range j {
			warnDisk := make(map[string]float64, 0)
			if v > config.GetConfig().AlterLimit {
				warnDisk[k] = v
				warnDiskList = append(warnDiskList, warnDisk)
			}
		}
	}
	return warnDiskList
}

func MonitorDisk() {
	for {
		warnDiskList := checkDisk()
		if len(warnDiskList) != 0 {
			newlog.Mylog("系统报警").Warn("磁盘报警了!!!")
			msg := "磁盘报警:已经超过警告值" + " " + FloatToString(config.GetConfig().AlterLimit)
			SendDingMsg(msg)
		}
		time.Sleep(time.Duration(config.GetConfig().Interval) * time.Second)
	}
}

//检查网络流量
func checkNet() {
	info, _ := net.IOCounters(true)
	recvMapList := make([]map[string]int, 0)
	sendMapList := make([]map[string]int, 0)
	for _, v := range info {
		recvMap := make(map[string]int)
		if int(v.BytesRecv)/1048576/100 > config.GetConfig().Recv {
			recvMap[v.Name] = int(v.BytesRecv) / 1048576
			recvMapList = append(recvMapList, recvMap)
		}
		if int(v.BytesSent)/1048576/100 > config.GetConfig().Send {
			sendMap := make(map[string]int)
			sendMap[v.Name] = int(v.BytesSent) / 1048576
			sendMapList = append(sendMapList, sendMap)
		}
		fmt.Println(v.BytesRecv/1048576, v.BytesSent/1048576)
	}
	recvResultJson, _ := json.Marshal(recvMapList)
	recvResult := string([]byte(recvResultJson))
	sendResultJson, _ := json.Marshal(sendMapList)
	sendResult := string([]byte(sendResultJson))
	if len(recvMapList) != 0 {
		newlog.Mylog("系统服务").Warn("进入流量过高，超过报警值: " + strconv.Itoa(config.GetConfig().Recv) + "M" + "当前流量值为:" + recvResult)
		SendDingMsg("进入流量过高，超过报警值: " + strconv.Itoa(config.GetConfig().Recv) + "M" + "当前流量值为:" + recvResult)
	}
	if len(sendMapList) != 0 {
		newlog.Mylog("系统服务").Warn("流出流量过高，超过报警值: " + strconv.Itoa(config.GetConfig().Send) + "M" + "当前流量值为:" + sendResult)
		SendDingMsg("流出流量过高，超过报警值: " + strconv.Itoa(config.GetConfig().Send) + "M" + "当前流量值为:" + sendResult)
	}
}
func MonitorNet() {
	for {
		checkNet()
		time.Sleep(time.Duration(config.GetConfig().Interval) * time.Second)
		//SendDingMsg("流量超过了警告值")
	}

}
