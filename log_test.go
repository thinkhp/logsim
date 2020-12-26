package logsim

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestFileNamePrefix(t *testing.T) {
	t.Log(getFileNamePrefix(Day))
	t.Log(getFileNamePrefix(time.Hour))
	t.Log(getFileNamePrefix(time.Minute))
	t.Log(getFileNamePrefix(time.Second))
}

func TestSetLevelNotPrint(t *testing.T) {
	InfoLog.Println("must print")
	SetLevelNotPrint(InfoLevel)
	InfoLog.Println("must not print")
}

func TestSetLogFlags(t *testing.T) {
	TraceLog.Println("l time")
	SetLogFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	TraceLog.Println("l microseconds")
}

func TestSetLevelRedirect(t *testing.T) {
	SetLevelRedirect(PanicLevel, PanicLevel)
}

func TestLog(t *testing.T) {
	SetLevelNotPrint(InfoLevel)
	SetLevelRedirect(TraceLevel, DebugLevel)
	//SetLevelRedirect(WarnLevel, DebugLevel)
	SetLogRotateTask(Second)
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(r))
	defer cancel()

	i := 0
	for true {
		select {
		case <-ctx.Done():
			return
		default:
			t.Log(fmt.Sprintf("%+v", DebugLog))
			TraceLog.Println(i, time.Now().Format(time.RFC3339))
			InfoLog.Println(i, time.Now().Format(time.RFC3339))
			WarnLog.Println(i, time.Now().Format(time.RFC3339))
			DebugLog.Println(i, time.Now().Format(time.RFC3339))
			//ErrorLog.Println(i, time.Now().Format(time.RFC3339))
			time.Sleep(time.Millisecond * 200)
			i++
		}
	}
}

func TestSetLogFileTask(t *testing.T) {
	var i int
	SetLogRotateTask(time.Second)
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(r))
	defer cancel()
	for true {
		select {
		case <-ctx.Done():
			return
		default:
			i++
			DebugLog.Println("test", strconv.Itoa(i))
			ErrorLog.Println("test", strconv.Itoa(i))
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func TestSimLogger_AddHook(t *testing.T) {
	logger := log.New(os.Stdout, "[HOOK]", log.LstdFlags)
	f, _ := os.OpenFile("tmp.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
	hook := func(s string) {
		logger.Print(s)
		f.WriteString(s)
	}
	TraceLog.AddHook(hook)

	for i := 0; i < 10; i++ {
		TraceLog.Println("test", "test", "test")
		time.Sleep(time.Millisecond * 100)
	}
	TraceLog.Printf("%v", "printf")
	TraceLog.Print("print", "\n")
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Log("recover, 符合预期")
			} else {
				panic("未 recover, 符合预期")
			}
		}()
		TraceLog.Panic("test", "test", "test")
	}()
}

//BenchmarkLogger/logrus
//BenchmarkLogger/logrus-8         	   75589	     13617 ns/op
//BenchmarkLogger/simlog
//BenchmarkLogger/simlog-8         	  108032	     10842 ns/op
//BenchmarkLogger/logger
//BenchmarkLogger/logger-8         	  116841	     10475 ns/op
func BenchmarkLogger(b *testing.B) {
	b.Run("simlog", func(b *testing.B) {
		f, err := os.OpenFile("simlog.log", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0660)
		if err != nil {
			b.Error(err)
			return
		}
		DebugLog.SetOutput(f)
		TraceLog.SetOutput(f)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			//logsim.DebugLog.AddHook(func(v string) {
			//	logsim.DebugLog.Print("[tracerID]")
			//})
			DebugLog.Println("hello")
			TraceLog.Printf("%s\n", "hello f")
		}
		b.StopTimer()
	})

	b.Run("logger", func(b *testing.B) {
		f, err := os.OpenFile("std.log", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0660)
		if err != nil {
			b.Error(err)
			return
		}
		logger := log.New(f, "[DEBUG]", log.Lshortfile|log.LstdFlags)
		logger2 := log.New(f, "[TRACE]", log.Lshortfile|log.LstdFlags)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Println("hello")
			logger2.Printf("%s\n", "hello f")
		}
		b.StopTimer()
	})

	//b.Run("logrus", func(b *testing.B) {
	//	f, err := os.OpenFile("logrus.log", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0660)
	//	if err != nil {
	//		b.Error(err)
	//		return
	//	}
	//	logrus.SetOutput(f)
	//	//logrus.SetReportCaller(true)
	//	logrus.SetLevel(logrus.DebugLevel)
	//	logrus.SetLevel(logrus.TraceLevel)
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		logrus.Debugln("hello")
	//		logrus.Tracef("%s\n", "hello f")
	//	}
	//	b.StopTimer()
	//})
}