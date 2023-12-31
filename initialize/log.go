package initialize

// 初始化日志
import "log"

// 未了快速实现demo，此处直接采用了官方log库，后面优化最好修改为zap日志库（我了解到这个库速度和使用量还是比较多的）

type MyLog struct {
	*log.Logger
	level int
}

var Log *MyLog

const (
	Debug int = iota
	Info
	Error
)

func NewLog() {
	Log = &MyLog{Logger: log.Default()}
}

func (l *MyLog) SetLevel(level int) {
	l.level = level
}

func (l *MyLog) Info(msgArr ...any) {
	l.print(Info, msgArr...)
}

func (l *MyLog) Debug(msgArr ...any) {
	l.print(Debug, msgArr...)
}

func (l *MyLog) Error(msgArr ...any) {
	l.print(Error, msgArr...)
}

func (l *MyLog) print(level int, msgArr ...any) {
	if l.level <= level {
		l.Logger.Print(msgArr...)
	}
}
