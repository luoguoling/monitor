//检查监控的文件是否有变化
package monitorlib

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mwmonitor/config"
	newlog "mwmonitor/logger"
	"mwmonitor/publib"
	"os"
	"strings"
	"time"
)

//获取文件md5值
func GetFileMd5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("os Open error")
		return "", err
	}
	md5 := md5.New()
	_, err = io.Copy(md5, file)
	if err != nil {
		fmt.Println("io copy error")
		return "", err
	}
	md5Str := hex.EncodeToString(md5.Sum(nil))
	return md5Str, nil
}

//判断文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//将md5写入文件
func WriteFile(s string) {
	var filename = ".filesmd5.txt"
	var f *os.File
	if checkFileIsExist(filename) {
		fmt.Println("存在这个文件")
		f, _ = os.OpenFile(filename, os.O_APPEND, 0666)
	} else {
		fmt.Println("不存在这个文件!!!")
		f, _ = os.Create(filename)
	}
	io.WriteString(f, s)
}

//获取上次这个文件的md5值
func ReadFileJson(filename string) string {
	fileopen, err := ioutil.ReadFile(".filesmd5.txt")
	if err != nil {
		panic(err)
	}
	filemap := make(map[string]string)
	err = json.Unmarshal([]byte(string(fileopen)), &filemap)
	if err != nil {
		panic(err)
	}
	return filemap[filename]
}

//获取上次文件的map
func GetOldMap() map[string]string {
	fileopen, err := ioutil.ReadFile(".filesmd5.txt")
	if err != nil {
		panic(err)
	}
	filemap := make(map[string]string)
	err = json.Unmarshal([]byte(string(fileopen)), &filemap)
	if err != nil {
		panic(err)
	}
	return filemap
}

//监控文件是否发生变化
func Notifyfiles(files []string) {
	var md5file = ".filesmd5.txt"
	md5map := make(map[string]string)
	changefile := make([]string, 0)
	for _, f := range files {
		md5sum, _ := GetFileMd5(f)
		md5map[f] = md5sum
	}
	if !checkFileIsExist(md5file) {
		md5mapjson, _ := json.Marshal(md5map)
		WriteFile(string(md5mapjson))
	}
	oldfilemap := GetOldMap()
	if len(oldfilemap) != len(md5map) {
		os.Remove(md5file)
		md5mapjson, _ := json.Marshal(md5map)
		WriteFile(string(md5mapjson))
	}
	for filename, _ := range md5map {
		if md5map[filename] != ReadFileJson(filename) {
			changefile = append(changefile, filename)
		}
	}
	if len(changefile) != 0 {
		currentime := publib.GetTime()
		changefilestr, err := json.Marshal(changefile)
		if err != nil {
			panic(err)
		}
		//result := string([]byte(changefilestr))
		result := string(changefilestr)
		sendstr := config.GetConfig().ProjectName + string(result) + "文件发生变动，请及时查看" + string(currentime)
		sendstrMsg := strings.Replace(sendstr, "\"", "", -1)
		//fmt.Println(sendstrMsg)
		SendDingMsg(sendstrMsg)
		newlog.Mylog("文件安全监控").Error(sendstrMsg)

	}

}
func Monitorfile() {
	files := config.GetConfig().Files
	for {
		Notifyfiles(files)
		time.Sleep(time.Duration(config.GetConfig().Interval) * time.Second)
	}
}
