package main

import (
	"flag"
	"fmt"
	_ "github.com/TomatoMr/visor/config"
	"io/ioutil"
	"log"
	"mwmonitor/config"
	"mwmonitor/monitorlib"
	_ "mwmonitor/monitorlib"
	"mwmonitor/publib"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

func main() {
	pool := make(chan struct{}, 10) //控制并发数量
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		fmt.Println("6060")
	}()
	var configPath string
	var start bool
	var stop bool
	var daemon bool
	var restart bool
	flag.StringVar(&configPath, "config", "config/config.yaml", "assign your config file: -config=your_config_file_path.")
	flag.BoolVar(&start, "start", false, "up your app, just like this: -start or -start=true|false.")
	flag.BoolVar(&stop, "stop", false, "down your app, just like this: -stop or -stop=true|false.")
	flag.BoolVar(&restart, "restart", false, "restart your app, just like this: -restart or -restart=true|false.")
	flag.BoolVar(&daemon, "d", false, "daemon, just like this: -start -d or -d=true|false.")
	flag.Parse()
	if err := config.InitConfig(configPath); err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	if start {
		if daemon {
			cmd := exec.Command(os.Args[0], "-start", "-config="+configPath)
			cmd.Start()
			os.Exit(0)
		}
		wg.Add(1)
		pool <- struct{}{} //控制并发数量
		fmt.Println("start.")
		Start()
		defer func() {
			<-pool
		}()
		wg.Wait()
	}

	if stop {
		Stop()
	}

	if restart {
		Restart()
	}

	//处理信号
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case <-sigs:
		return
	}

}

func Start() {
	defer wg.Done()
	ioutil.WriteFile(config.GetConfig().Pid, []byte(fmt.Sprintf("%d", os.Getpid())), 0666) //记录pid
	monitoritems := config.GetConfig().MonitorItems
	if publib.IsValueInList("MonitorCpu", monitoritems) {
		go monitorlib.MonitorCpu()
	}
	if publib.IsValueInList("MonitorMem", monitoritems) {
		go monitorlib.MonitorMem()
	}
	if publib.IsValueInList("MonitorExceptionLog", monitoritems) {
		go monitorlib.MonitorExceptionLog()
	}
	if publib.IsValueInList("Monitorfile", monitoritems) {
		go monitorlib.Monitorfile()
	}
	if publib.IsValueInList("MonitorAllProRunning", monitoritems) {
		go monitorlib.MonitorAllProRunning()
	}
	if publib.IsValueInList("MonitorIptable", monitoritems) {
		go monitorlib.MonitorIptable()
	}
	if publib.IsValueInList("MonitorDisk", monitoritems) {
		go monitorlib.MonitorDisk()
	}
	if publib.IsValueInList("MonitorLoad", monitoritems) {
		go monitorlib.MonitorLoad()
	}

	//go monitorlib.MonitorNet()
	if publib.IsValueInList("MonitorDangerCmdLog", monitoritems) {
		go monitorlib.MonitorDangerCmdLog()
	}
	if publib.IsValueInList("MonitorCheckIp", monitoritems) {
		go monitorlib.MonitorCheckIp()
	}

	//fmt.Println("开始执行任务计划")
	//go monitorlib.CheckCron()
	//fmt.Println("任务计划执行完毕")
	//time.Sleep(time.Duration(config.GetConfig().Interval) *2* time.Second)
}

func Stop() {
	pid, _ := ioutil.ReadFile(config.GetConfig().Pid)
	cmd := exec.Command("kill", "-9", string(pid))
	cmd.Start()
	fmt.Println("kill ", string(pid))
	os.Remove(config.GetConfig().Pid) //清除pid
	os.Exit(0)
}

func Restart() {
	fmt.Println("restarting...")
	pid, _ := ioutil.ReadFile(config.GetConfig().Pid)
	stop := exec.Command("kill", "-9", string(pid))
	stop.Start()
	fmt.Println(os.Args[0])
	start := exec.Command(os.Args[0], "-start", "-d")
	start.Start()
	os.Exit(0)
	//Stop()
	//time.Sleep(time.Duration(4) * time.Second)
	//Start()
}

//monitorlib.MonitorCpu()
//monitorlib.MonitorExceptionLog()

//monitorlib.SendDingMsg("awerwerwtest")
//files := config.GetConfig().Files
//fmt.Println(files)
//monitorlib.Notifyfiles(files)
