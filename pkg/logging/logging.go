package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	F             *os.File
	DefaultPrefix = "INFO"
	logPrefix     = ""
	logger        *log.Logger
	levelFlags    = []string{"TRACE", "INFO", "WARN", "ERROR", "FATAL"}
)

type Level int

const (
	TRACE Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func SetUp() {
	// 直接使用标准输出，不需要文件
	logger = log.New(os.Stdout, DefaultPrefix, log.LstdFlags|log.Lshortfile)
}

func Trace(v ...interface{}) {
	setPrefix(TRACE)
	logger.Println(v)
}
func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}
func Warn(v ...interface{}) {
	setPrefix(WARN)
	logger.Println(v)
}
func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}
func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Println(v)
}

func setPrefix(level Level) {
	_, filename, line, ok := runtime.Caller(0)
	if ok {
		logPrefix = fmt.Sprintf("【%s】-【%s:%d】-", levelFlags[level], filename, line)
	} else {
		logPrefix = fmt.Sprintf("【%s】", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}
