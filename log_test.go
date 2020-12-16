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
