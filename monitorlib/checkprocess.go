//检查监控的进程是否正常运行
package monitorlib

import (
	"encoding/json"
	"mwmonitor/config"
	newlog "mwmonitor/logger"
	"mwmonitor/publib"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

//根据进程名判断进程是否运行
func runInWindows(cmd string) (string, error) {
	result, err := exec.Command("cmd", "/c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}
func runInLinux(cmd string) (string, error) {
	//fmt.Println("Running Linux cmd:" + cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

func RunCommand(cmd string) (string, error) {
	if runtime.GOOS == "windows" {
		return runInWindows(cmd)
	} else {
		return runInLinux(cmd)
	}
}
//根据进程名判断进程是否运行
func CheckProRunning(serverName string) (bool, error) {
	a := `ps aux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	if err != nil {
		return false, err
	}
	return pid != "", nil
}
//检查所有进程是否在运行
func CheckAllProRunning()  {
	currentime := publib.GetTime()
	notRunningPros := make([]string,0)
	processes := config.GetConfig().Processes
	for _,process := range processes{
		isRunning,_ := CheckProRunning(process)
		if isRunning == false{
			notRunningPros = append(notRunningPros,process)
		}
	}
	notRunningPosString,_ := json.Marshal(notRunningPros)
	result := string([]byte(notRunningPosString))
	pubip,err := publib.GetPubIp()
	if err != nil{
		newlog.Mylog("程序自身错误").Error("获取公有ip报错!!!")
	}
	sendstr := config.GetConfig().ProjectName + ":::" + pubip + ":::" + result + "进程没有启动!!!" + string(currentime)
	sendstrMsg := strings.Replace(sendstr,"\"","",-1)
	if len(notRunningPros)  != 0{
		newlog.Mylog("进程监控").Warn(sendstrMsg)
		SendDingMsg(sendstrMsg)
	}
}
func MonitorAllProRunning(){
	for{
		CheckAllProRunning()
		time.Sleep(time.Duration(config.GetConfig().Interval) * time.Second)
	}
}
