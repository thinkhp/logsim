package logsim

// 1.输出方法与标准库一致
// 2.允许部分级别的日志不输出
// 3.允许部分级别的日志输出到相同文件中
// 4.日志切割,级别为:天,时,分钟,秒
// TODO 5.日志传送到server
// TODO 6.以 json 或者 text 的形式输出
// TODO 7.print() runtime.Call(1)

import (
	"io"
	"log"
	"os"
)

type level int8

const (
	TraceLevel level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

const (
	//defaultOutput = os.Stdout
	defaultLevel   = TraceLevel
	defaultFlags   = log.LstdFlags | log.Lshortfile
	defaultLogPath = "./log/"
)

type devNull int

var stdNull = devNull(0)

func (devNull) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (devNull) Close() error {
	return nil
}

type Logger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

type iLogger struct {
	levelPrint level
	levelFile  level
	*log.Logger
}

//定义logger, 传入参数 文件，前缀字符串，flag标记
func New(lv level, out io.Writer, prefix string, flag int) *iLogger {
	l := new(iLogger)
	l.levelPrint = lv
	l.levelFile = lv
	l.Logger = log.New(out, prefix, flag)

	return l
}

var cfgLogPath = defaultLogPath

func SetLogPath(dir string) {
	cfgLogPath = dir
	// 如果path指定了一个已经存在的目录，MkdirAll不做任何操作并返回nil。
	err := os.MkdirAll(cfgLogPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

var cfgFlags = defaultFlags
var cfgLevel = defaultLevel

// ls 级别的日志即使调用,也不会输出
func SetLevelNotPrint(ls ...level) {
	for _, v := range ls {
		switch v {
		case TraceLevel:
			TraceLog.(*iLogger).SetOutput(stdNull)
		case PanicLevel:
		case FatalLevel:
		case ErrorLevel:
			ErrorLog.(*iLogger).SetOutput(stdNull)
		case WarnLevel:
			WarnLog.(*iLogger).SetOutput(stdNull)
		case InfoLevel:
			InfoLog.(*iLogger).SetOutput(stdNull)
		case DebugLevel:
			DebugLog.(*iLogger).SetOutput(stdNull)
		}
	}
}

// 将 src 级别的日志输出到 redirect 级别的日志文件
func SetLevelRedirect(src, redirect level) {
	allLog[src].(*iLogger).levelFile = redirect
}

var allLevel = []level{
	//PanicLevel,
	//FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}
var allLog map[level]Logger

var TraceLog Logger
var DebugLog Logger
var InfoLog Logger
var WarnLog Logger
var ErrorLog Logger

func init() {
	TraceLog = New(TraceLevel, os.Stdout, "[TRACE] ", cfgFlags)
	DebugLog = New(DebugLevel, os.Stdout, "[DEBUG] ", cfgFlags)
	InfoLog = New(InfoLevel, os.Stdout, "[INFO ] ", cfgFlags)
	WarnLog = New(WarnLevel, os.Stdout, "[WARN ] ", cfgFlags)
	ErrorLog = New(ErrorLevel, os.Stderr, "[ERROR] ", cfgFlags)

	allLog = map[level]Logger{
		ErrorLevel: ErrorLog,
		WarnLevel:  WarnLog,
		InfoLevel:  InfoLog,
		DebugLevel: DebugLog,
		TraceLevel: TraceLog,
	}

	SetLevelRedirect(TraceLevel, DebugLevel)
	SetLevelRedirect(InfoLevel, DebugLevel)
	SetLevelRedirect(WarnLevel, DebugLevel)
}

// 在Linux中,因为不存在文件保护,所以要检查文件名所对应的指针是否存在
func checkFile(filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
}
