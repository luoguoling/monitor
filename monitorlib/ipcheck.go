//检查当前登录ip是不是要求登录的ip
package monitorlib

import (
	"fmt"
	"mwmonitor/config"
	"mwmonitor/logger"
	"time"
)

func CheckIp() bool {
	whiteIp := config.GetConfig().WriteIp
	currentip, err := runInLinux("echo $SSH_CLIENT |awk ' { print $1 }'")
	fmt.Println(currentip)
	if err != nil {
		fmt.Println("获取当前ip错误!!!")
	}
	for _, ip := range whiteIp {
		fmt.Println("ip is ...")
		fmt.Println(ip)
		if currentip == ip {
			return true
		}
	}
	return false
}
func MonitorCheckIp() {
	for {
		Ipsec := CheckIp()
		if Ipsec == false {
			logger.Mylog("系统安全报错").Error("当前ip非允许ip登录!!!!")
			SendDingMsg("当前ip非允许ip登录!!!")
		}
		time.Sleep(time.Duration(config.GetConfig().Interval) * time.Second)
	}

}
