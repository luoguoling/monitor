//记录危险命令操作
package monitorlib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mwmonitor/config"
	newlog "mwmonitor/logger"
	"mwmonitor/publib"
	"path/filepath"
	"strings"
	"time"
)

func listDir(path string, indent int) (files []string) {

	dir, err := filepath.Abs(path)

	if err != nil {
		fmt.Printf("err : %s \n", err)
		return
	}

	//fmt.Printf("%sDir: %s\n", strings.Repeat(" ", indent*4), dir)
	finfos, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("err : %s \n", err)
		return
	}
	//dangerwords := config.GetConfig().DangerWords
	// 遍历子目录或者文件
	for _, fi := range finfos {
		_ = fi
		// 如果是目录，则递归输出
		//if fi.IsDir() {
		//	listDir(dir+string(os.PathSeparator)+fi.Name(), indent+1)
		//	continue
		//}
		//dangerwords := config.GetConfig().DangerWords
		// 如果是文件，则直接输出文件名
		//fmt.Printf("%s%s\n", "/", fi.Name())
		file := fmt.Sprintf("%s%s%s%s", dir, "/", fi.Name(), "/cmd.txt")
		files = append(files, file)

	}
	return files

}

func CheckDangerCmd() {
	currentime := publib.GetTime()
	var sendstr string
	var dangerslices []map[string]string
	dangerwords := config.GetConfig().DangerWords
	files := listDir("/usr/bin/.hist", 0)
	//fmt.Println(files)
	for _, file := range files {
		kTmp := file
		//fmt.Println("危险命令字典")
		//fmt.Println(dangerwords)
		for _, exceptionkeyword := range dangerwords {
			//fmt.Println("执行到错误关键字检查")
			if CheckFileContainsStr(file, exceptionkeyword) {
				vTmp := exceptionkeyword
				var exceptionmap = make(map[string]string)
				exceptionmap[kTmp] = vTmp
				dangerslices = append(dangerslices, exceptionmap)
			}

		}
	}
	dangerlicesstr, err := json.Marshal(dangerslices)
	//fmt.Println(dangerlicesstr)
	if err != nil {
		newlog.MyLogger.Warn("解析异常")
	}
	result := string([]byte(dangerlicesstr))
	//fmt.Println(result)
	//fmt.Println("exceptionslics的长度是:",len(exceptionslices))
	//fmt.Println(exceptionslices)
	sendstr = config.GetConfig().ProjectName + " " + string(result) + "文件含有危险命令信息!!!" + string(currentime)
	sendstrMsg := strings.Replace(sendstr, "\"", "", -1)
	if len(dangerslices) != 0 {
		newlog.Mylog("文件含有危险命令信息").Warn(sendstrMsg)
		//fmt.Println("钉钉发送报警消息:", sendstrMsg)
		SendDingMsg(sendstrMsg)
	}

}

func MonitorDangerCmdLog() {
	for {
		CheckDangerCmd()
		CheckCron()
		time.Sleep(time.Duration(config.GetConfig().Interval) * time.Second)
	}
}
