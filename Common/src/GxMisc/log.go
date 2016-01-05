/**
作者:guangbo
模块：日志模块
说明：
创建时间：2015-10-30
**/
package GxMisc

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	TraceLevel = iota
	DebugLevel
	WarnLevel
	ErrorLevel
	InfoLevel
	FatalLevel
)

var logger *log.Logger
var level = TraceLevel
var logFile *os.File
var isOutputScreen = true

// 获取日志级别
func GetLevel() int {
	return level
}

func SetLevel(l int) {
	if l > FatalLevel || l < TraceLevel {
		level = TraceLevel
	} else {
		level = l
	}
}

func SetIsOutputScreen(isOutput bool) {
	isOutputScreen = isOutput
}

// Initialization logger
func InitLogger(logFileName string) {
	var err error
	pid := os.Getpid()                                                                 //Get system pid.
	pidStr := strconv.FormatInt(int64(pid), 10)                                        //Convert int to string.
	logFileName = "log/" + logFileName + "_" + pidStr + ".log"                         //Get config file path.
	logFile, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666) //Open config file.

	//Return error.
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//Create new log.
	logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Println("log to file sample")
}

func Trace(format string, v ...interface{}) {
	if level <= TraceLevel {
		var str string
		str = "[T] " + format
		str = fmt.Sprintf(str, v...)
		logger.Output(2, str)

		if isOutputScreen {
			fmt.Println(str)
		}
	}
}

func Debug(format string, v ...interface{}) {
	if level <= DebugLevel {
		var str string
		str = "[D] " + format
		str = fmt.Sprintf(str, v...)
		logger.Output(2, str)

		if isOutputScreen {
			fmt.Println(str)
		}
	}
}

func Warn(format string, v ...interface{}) {
	if level <= WarnLevel {
		var str string
		str = "[W] " + format
		str = fmt.Sprintf(str, v...)
		logger.Output(2, str)

		if isOutputScreen {
			fmt.Println(str)
		}
	}
}

func Error(format string, v ...interface{}) {
	if level <= ErrorLevel {
		var str string
		str = "[E] " + format
		str = fmt.Sprintf(str, v...)
		logger.Output(2, str)

		if isOutputScreen {
			fmt.Println(str)
		}
	}
}

func Info(format string, v ...interface{}) {
	if level <= InfoLevel {
		var str string
		str = "[I] " + format
		str = fmt.Sprintf(str, v...)
		logger.Output(2, str)

		if isOutputScreen {
			fmt.Println(str)
		}
	}
}

func Fatal(format string, v ...interface{}) {
	if level <= FatalLevel {
		var str string
		str = "[F] " + format
		str = fmt.Sprintf(str, v...)
		logger.Output(2, str)

		if isOutputScreen {
			fmt.Println(str)
		}
	}
}
