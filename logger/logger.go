package logger

import "fmt"

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

type SimpleLogger struct {
	level LogLevel
}

func (sl *SimpleLogger) Debug(v ...interface{}) {
	if sl.level <= LogLevelDebug {
		fmt.Print("[Debug] ")
		fmt.Println(v...)
	}
}
func (sl *SimpleLogger) Debugf(format string, v ...interface{}) {
	if sl.level <= LogLevelDebug {
		fmt.Printf("[Debug] "+format+"\n", v...)
	}
}
func (sl *SimpleLogger) Info(v ...interface{}) {
	if sl.level <= LogLevelInfo {
		fmt.Print("[Info] ")
		fmt.Println(v...)
	}
}
func (sl *SimpleLogger) Infof(format string, v ...interface{}) {
	if sl.level <= LogLevelInfo {
		fmt.Printf("[Info] "+format+"\n", v...)
	}
}
func (sl *SimpleLogger) Warn(v ...interface{}) {
	if sl.level <= LogLevelWarn {
		fmt.Print("[Warn] ")
		fmt.Println(v...)
	}
}
func (sl *SimpleLogger) Warnf(format string, v ...interface{}) {
	if sl.level <= LogLevelWarn {
		fmt.Printf("[Warn] "+format+"\n", v...)
	}
}
func (sl *SimpleLogger) Error(v ...interface{}) {
	if sl.level <= LogLevelError {
		fmt.Print("[Error] ")
		fmt.Println(v...)
	}
}
func (sl *SimpleLogger) Errorf(format string, v ...interface{}) {
	if sl.level <= LogLevelError {
		fmt.Printf("[Error] "+format+"\n", v...)
	}
}
func (sl *SimpleLogger) Fatal(v ...interface{}) {
	if sl.level <= LogLevelFatal {
		fmt.Print("[Fatal] ")
		fmt.Println(v...)
	}
}
func (sl *SimpleLogger) Fatalf(format string, v ...interface{}) {
	if sl.level <= LogLevelFatal {
		fmt.Printf("[Fatal] "+format+"\n", v)
	}
}

var simpleLogger = SimpleLogger{LogLevelInfo}

func Debug(v ...interface{}) {
	simpleLogger.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	simpleLogger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	simpleLogger.Info(v...)
}
func Infof(format string, v ...interface{}) {
	simpleLogger.Infof(format, v...)
}

func Warn(v ...interface{}) {
	simpleLogger.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	simpleLogger.Warnf(format, v...)
}

func Error(v ...interface{}) {
	simpleLogger.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	simpleLogger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	simpleLogger.Fatal(v...)
}
func Fatalf(format string, v ...interface{}) {
	simpleLogger.Fatalf(format, v...)
}
