//定义公共方法
package publib

import (
	"fmt"
	"io/ioutil"
	"log"
	"mwmonitor/config"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

//获取公有ip
func GetPubIp1() (string, error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		//logger.Mylog("程序自身错误").Error("程序报错崩溃退出!!!")
		os.Stderr.WriteString("\n")
		//os.Exit(1)
		return "ipnull", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return string(data), nil
}

//获取公有ip2
func GetPubIp() (string, error) {
	resp, err := http.Get("http://myip.ipip.net")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		//logger.Mylog("程序自身错误").Error("程序报错崩溃退出!!!")
		os.Stderr.WriteString("\n")
		//os.Exit(1)
		return "ipnull", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	stringdata := string(data)
	ip1 := strings.Split(stringdata, " ")[1]
	ip := strings.Split(ip1, ":")[0]
	fmt.Println(ip)
	return ip, nil
}

//获取当前时间
func GetTime() string {
	currentime := time.Now().Format("2006-01-02 15:04:05")
	return currentime

}

//判断目录是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//获取监控项目名字
func GetProjectName() string {
	projectname := config.GetConfig().ProjectName
	return projectname
}
func Typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

//获取本地ip
func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}

// 获取局域网ip地址
func GetLocaHonst() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}

	}
	return ""
}

//获取内部ip
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	return localAddr.IP.String()
}

//判断slice里面是否含有某个字符串
func IsValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

//处理panic,用try catch形式
func Try(fun func(),handler func(interface{}))  {
	defer func() {
		if err := recover();err!=nil{
			handler(err)
		}
	}()
	fun()

}
