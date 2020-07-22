package main

import (
	"flag"
	"fmt"
	_ "github.com/TomatoMr/visor/config"
	"io/ioutil"
	"mwmonitor/config"
	"mwmonitor/monitorlib"
	_ "mwmonitor/monitorlib"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)
var wg sync.WaitGroup



func main() {
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
			fmt.Println("start.")
			Start()
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
	go monitorlib.MonitorCpu()
	go monitorlib.MonitorMem()
	go monitorlib.MonitorExceptionLog()
	go monitorlib.Monitorfile()
	go monitorlib.MonitorAllProRunning()
	go monitorlib.MonitorIptable()
	go monitorlib.MonitorDisk()
	go monitorlib.MonitorLoad()
	//go monitorlib.MonitorNet()
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

