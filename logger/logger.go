package logger

import "fmt"

var simpleLogger = SimpleLogger{LogLevelInfo}

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
	LogLevelTrace LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

type SimpleLogger struct {
	level LogLevel
}

func (sl *SimpleLogger) Trace(v ...interface{}) {
	if sl.level <= LogLevelTrace {
		arr := []interface{}{"[Trace]"}
		for _, item := range v {
			arr = append(arr, item)
		}

		fmt.Println(arr...)
	}
}
func (sl *SimpleLogger) Tracef(format string, v ...interface{}) {
	if sl.level <= LogLevelTrace {
		fmt.Printf("[Trace] "+format+"\n", v...)
	}
}

func (sl *SimpleLogger) Debug(v ...interface{}) {
	if sl.level <= LogLevelDebug {
		arr := []interface{}{"[Debug]"}
		for _, item := range v {
			arr = append(arr, item)
		}

		fmt.Println(arr...)
	}
}
func (sl *SimpleLogger) Debugf(format string, v ...interface{}) {
	if sl.level <= LogLevelDebug {
		fmt.Printf("[Debug] "+format+"\n", v...)
	}
}
func (sl *SimpleLogger) Info(v ...interface{}) {
	if sl.level <= LogLevelInfo {
		arr := []interface{}{"[Info ]"}
		for _, item := range v {
			arr = append(arr, item)
		}
		fmt.Println(arr...)
	}
}
func (sl *SimpleLogger) Infof(format string, v ...interface{}) {
	if sl.level <= LogLevelInfo {
		fmt.Printf("[Info ] "+format+"\n", v...)
	}
}
func (sl *SimpleLogger) Warn(v ...interface{}) {
	if sl.level <= LogLevelWarn {
		arr := []interface{}{"[Warn ]"}
		for _, item := range v {
			arr = append(arr, item)
		}
		fmt.Println(arr...)
	}
}
func (sl *SimpleLogger) Warnf(format string, v ...interface{}) {
	if sl.level <= LogLevelWarn {
		fmt.Printf("[Warn ] "+format+"\n", v...)
	}
}
func (sl *SimpleLogger) Error(v ...interface{}) {
	if sl.level <= LogLevelError {
		arr := []interface{}{"[Error]"}
		for _, item := range v {
			arr = append(arr, item)
		}
		fmt.Println(arr...)
	}
}
func (sl *SimpleLogger) Errorf(format string, v ...interface{}) {
	if sl.level <= LogLevelError {
		fmt.Printf("[Error] "+format+"\n", v...)
	}
}
func (sl *SimpleLogger) Fatal(v ...interface{}) {
	if sl.level <= LogLevelFatal {
		arr := []interface{}{"[Fatal]"}
		for _, item := range v {
			arr = append(arr, item)
		}
		fmt.Println(arr...)

	}
}
func (sl *SimpleLogger) Fatalf(format string, v ...interface{}) {
	if sl.level <= LogLevelFatal {
		fmt.Printf("[Fatal] "+format+"\n", v)
	}
}

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
