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

func Debug(v ...interface{}) {
	fmt.Print("[Debug] ")
	fmt.Println(v...)
}
func Debugf(format string, v ...interface{}) {
	fmt.Printf("[Debug] "+format+"\n", v...)
}
func Info(v ...interface{}) {
	fmt.Print("[Info] ")
	fmt.Println(v...)
}
func Infof(format string, v ...interface{}) {
	fmt.Printf("[Info] "+format+"\n", v...)
}
func Warn(v ...interface{}) {
	fmt.Print("[Warn] ")
	fmt.Println(v...)
}
func Warnf(format string, v ...interface{}) {
	fmt.Printf("[Warn] "+format+"\n", v...)
}
func Error(v ...interface{}) {
	fmt.Print("[Error] ")
	fmt.Println(v...)
}
func Errorf(format string, v ...interface{}) {
	fmt.Printf("[Error] "+format+"\n", v...)
}
func Fatal(v ...interface{}) {
	fmt.Print("[Fatal] ")
	fmt.Println(v...)
}
func Fatalf(format string, v ...interface{}) {
	fmt.Printf("[Fatal] "+format+"\n", v)
}
