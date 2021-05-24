//报警方法后面可以增加其他报警方法
package monitorlib

import (
	"fmt"
	"mwmonitor/config"
	newlog "mwmonitor/logger"
	"mwmonitor/publib"
	"net/http"
	"strings"
)

//钉钉报警
func SendDingMsg(msg string) {
	defer func() {
		if err := recover();err!=nil{
			newlog.Mylog("程序日志").Error("发送报警异常")
		}
	}()
	//webHook := `https://oapi.dingtalk.com/robot/send?access_token=187918ed0afc579ee5e458f0bf23c84a1bafdd1782b683ad42873b4d41bba0d7`
	webHook := config.GetConfig().WebHook
	pubip, _ := publib.GetPubIp()
	msg = config.GetConfig().ProjectName + "  服务器ip:  " + pubip + " " + msg
	content := `{"msgtype": "text",
		"text": {"content": "` + msg + `"}
	}`
	//创建一个请求
	req, err := http.NewRequest("POST", webHook, strings.NewReader(content))
	if err != nil {
		fmt.Println(err)
		fmt.Println("钉钉报警请求异常")
	}

	client := &http.Client{}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		fmt.Println(err)
		fmt.Println("顶顶报发送异常!!!")
	}
	defer resp.Body.Close()

	//logger.MyLogger.Error("aaerw")
	//logrus.WithFields(logrus.Fields{"aa":"aa","username":"rolin"}).Info("aaaa")
	//logrus.Error("aaaa")

}
