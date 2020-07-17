package logsim

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	fmt.Println(getFileNamePrefix(Day))
	fmt.Println(getFileNamePrefix(time.Hour))
	fmt.Println(getFileNamePrefix(time.Minute))
	fmt.Println(getFileNamePrefix(time.Second))
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
			fmt.Println(DebugLog)
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
	for true {
		i++
		DebugLog.Println(`test`, strconv.Itoa(i))
		time.Sleep(time.Millisecond * 100)
	}
}
