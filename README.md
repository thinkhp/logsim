#### 精简的日志



```
import "github.com/thinkhp/logsim"
```



**example**

```
import (
	"github.com/thinkhp/logsim"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	logsim.SetLevelNotPrint(logsim.InfoLevel) //stdout, stderr is null

	logsim.SetLevelRedirect(logsim.TraceLevel, logsim.DebugLevel) //trace file is debug.*log
	logsim.SetLevelRedirect(logsim.WarnLevel, logsim.DebugLevel) //level warn file is debug*.log

	logsim.SetLogRotateTask(logsim.Day) //stdout and stderr is file;rotate in days

	logsim.TraceLog.Println(time.Now().Format(time.RFC3339))
	logsim.InfoLog.Println(time.Now().Format(time.RFC3339))
	logsim.WarnLog.Println(time.Now().Format(time.RFC3339))
	logsim.DebugLog.Println(time.Now().Format(time.RFC3339))
	logsim.ErrorLog.Println(time.Now().Format(time.RFC3339))
}
```

