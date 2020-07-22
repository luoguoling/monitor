package monitorlib
import (
	"fmt"
	"mwmonitor/config"
	"mwmonitor/logger"
	"strings"
	"time"
)
func CheckIptable(){
	cmd := "iptables -t filter -nL INPUT"
	aa,err := runInLinux(cmd)
	if err != nil{
		logger.Mylog("程序自身错误").Error("程序执行linux防火墙命令报错!!!")
	}
	iptablenum := strings.Split(aa,"\n")
	fmt.Println(len(iptablenum))
	if len(iptablenum) < 5{
		logger.Mylog("应用报错").Error("防火墙都没有开启!!!!")
		SendDingMsg("防火墙没有开启!!!")
	}
}
func MonitorIptable()  {
	for{
		CheckIptable()
		time.Sleep(time.Duration(config.GetConfig().Interval) * time.Second)
	}

}