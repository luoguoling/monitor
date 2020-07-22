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
	"time"
)
//获取公有ip
func GetPubIp() (string,error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		//logger.Mylog("程序自身错误").Error("程序报错崩溃退出!!!")
		os.Stderr.WriteString("\n")
		//os.Exit(1)
		return "ipnull",err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return string(data),nil
}

//获取当前时间
func GetTime() string {
	currentime := time.Now().Format("2006-01-02 15:04:05")
	return currentime

}
//判断目录是否存在
func PathExists(path string) (bool) {
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
func GetProjectName() string{
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
//获取外部ip
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




