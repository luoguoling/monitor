//查看文件是否含有异常信息，及时报警写入日志es
package monitorlib
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"mwmonitor/config"
	newlog "mwmonitor/logger"
	"mwmonitor/publib"
	"os"
	"strings"
	"time"
)
//判断某个文件是否含有某一个字符串
func CheckFileContainsStr(path,str string) bool{
	fi, err := os.Open(path)
	exist := false
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if strings.Contains(string(a), str){
			exist = true
			return exist
		}
		if c == io.EOF {
			break
		}
	}
	return exist
}
func CheckExceptionLog()  {
	currentime := publib.GetTime()
	//errfiles := []string{"/root/a.txt","/root/b.txt"}
	//exceptionkeywords := []string{"err","warning", "exception"}
	errfiles := config.GetConfig().Errfiles
	exceptionkeywords := config.GetConfig().ExceptionKeywords
	//newlog.MyLogger.Warn("日志信息!!!!")
	var exceptionslices []map[string]string
	//var exceptionmap = make(map[string]string)
	var sendstr string
	for _,errfile := range errfiles{
		kTmp := errfile
		for _,exceptionkeyword := range exceptionkeywords{
			if CheckFileContainsStr(errfile,exceptionkeyword){
				vTmp := exceptionkeyword
				var exceptionmap = make(map[string]string)
				exceptionmap[kTmp] = vTmp
				exceptionslices = append(exceptionslices,exceptionmap)
			}
		}
		//fmt.Println(exceptionslices)

	}
	fmt.Println("exceptionslics的长度是111",exceptionslices)
	exceptionslicesstr,err := json.Marshal(exceptionslices)
	if err != nil{
		newlog.MyLogger.Warn("解析异常")
	}
	result := string([]byte(exceptionslicesstr))
	fmt.Println("exceptionslics的长度是:",len(exceptionslices))
	fmt.Println(exceptionslices)
	sendstr = config.GetConfig().ProjectName +" " + string(result) + "文件含有错误异常信息!!!" + string(currentime)
	sendstrMsg := strings.Replace(sendstr,"\"","",-1)
	if len(exceptionslices)  != 0{
		newlog.Mylog("文件报错").Warn(sendstrMsg)
		fmt.Println("钉钉发送报警消息:",sendstrMsg)
		SendDingMsg(sendstrMsg)
	}
}
func MonitorExceptionLog(){
	for{
		CheckExceptionLog()
		time.Sleep(time.Duration(config.GetConfig().Interval) * time.Second)
	}
}


