//查看任务计划是否有异常
package monitorlib

import (
	"fmt"
	"mwmonitor/config"
	"mwmonitor/logger"
	"strconv"
	"strings"
	"time"
)

func CheckCron() {
	fmt.Println("开始执行任务计划")
	cmd := "crontab -l|wc -l"
	aa, err := runInLinux(cmd)
	if err != nil {
		logger.Mylog("程序自身错误").Error("程序执行linux任务计划条数统计报错!!!")
	}
	cronnum := strings.Split(aa, "\n")
	normalcronnum := config.GetConfig().Cronnum
	crontnumint, _ := strconv.Atoi(cronnum[0])
	fmt.Println("判断任务计划条数")
	fmt.Println(crontnumint, normalcronnum)
	if crontnumint > normalcronnum {
		logger.Mylog("应用报错").Error(" 任务计划条数异常!!!!")
		SendDingMsg("任务计划条数异常!!!")
	}
}
func MonitorCron() {
	for {
		CheckCron()
		time.Sleep(time.Duration(config.GetConfig().Interval) * time.Second)
	}
}
