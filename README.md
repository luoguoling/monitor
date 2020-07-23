# monitor
服务器监控
1.系统监控：cpu 内存 负载 磁盘 io
2.安全监控:对重要文件的md5值进行监控
3.日志监控：监控文件种的错误日志并报警
4.进程监控：监控重要的进程是否运行
5.流量监控

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
