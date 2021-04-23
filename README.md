# monitor
主要监控项目：(可以自由定义监控选项，定义监控参数选择)

监控机器内存

监控指定文件是否含有异常日志

监控文件md5是否有变话

监控进程是否在运行

监控防火墙是否开启

监控磁盘利用率

监控服务器负载

监控操作日志是否含有定义的危险命令

监控当前登录ip是否是安全ip

监控任务计划运行条数是否有异常

使用方法：
nohup ./main -start -config config/config.yaml -d &
nohup ./main -stop
 
v1版本问题：cpu消耗有点儿高，已经搞定。
性能排查思路如下：
golang性能问题简单判断判断：
1.在main函数导入包_ "net/http/pprof"
2.在main函数添加以下语句
go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		fmt.Println("启用6060端口")
	}()
3.安装graphviz
4.启动以下命令:go tool pprof -http=0.0.0.0:6061 http://127.0.0.1:6060/debug/pprof/profile?seconds=30
5.通过图形界面查看消耗过高函数:http://ip:6061/ui/flamegraph
