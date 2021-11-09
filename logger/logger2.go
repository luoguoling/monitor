package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"mwmonitor/config"
	"mwmonitor/publib"
	"os"
)

//　初始化zaplogger日志库
var MyLogger *zap.Logger

func initLogger(logpath string) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logpath, // 日志文件路径
		MaxSize:    10,      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,      // 日志文件最多保存多少个备份
		MaxAge:     7,       // 文件最多保存多少天
		Compress:   true,    // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)
	//IoWriter := getWriter(logpath)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
	//core := zapcore.NewCore(
	//	zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
	//	zapcore.NewMultiWriteSyncer(zapcore.AddSync(IoWriter), zapcore.AddSync(&hook)), // 打印到控制台和文件
	//	atomicLevel,                                                                     // 日志级别
	//)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	pubip, _ := publib.GetPubIp()
	//if err != nil{
	//	Mylog("程序自身报错").Error("程序获取ip报错!!!")
	//}
	filed := zap.Fields(zap.String("ip", pubip))
	// 构造日志
	logger := zap.New(core, caller, development, filed)

	//logger.Info("log 初始化成功")
	return logger
}

func init() {
	//logPath := beego.AppConfig.String("LogPath")
	if !publib.PathExists("logs") {
		os.Mkdir("logs", 0666)
	}
	logPath := "logs/monitor.log"
	MyLogger = initLogger(logPath)
}
func Mylog(logtype string) *zap.Logger {
	projectname := config.GetConfig().ProjectName
	aa := MyLogger.With(zap.String("projectname", projectname), zap.String("logtype", logtype))
	return aa
}
