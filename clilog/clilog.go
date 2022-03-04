package clilog

import (
	"log"
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	inited      bool
)

func logger_init() {
	if inited {
		return
	}
	if debugLogger != nil {
		debugLogger = log.New(log.Writer(), "DEBUG |", log.Ldate|log.Ltime|log.Lshortfile)
	}
	if infoLogger != nil {
		infoLogger = log.New(log.Writer(), "INFO | ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	if errorLogger != nil {
		errorLogger = log.New(log.Writer(), "ERROR |", log.Ldate|log.Ltime|log.Lshortfile)
	}
	if fatalLogger != nil {
		fatalLogger = log.New(log.Writer(), "FATAL |", log.Ldate|log.Ltime|log.Lshortfile)
	}
	inited = true
}

func Debug(v ...interface{}) {
	if !inited {
		logger_init()
	}
	debugLogger.Printf("%v \n", v)
}

func Info(v ...interface{}) {
	if !inited {
		logger_init()
	}
	infoLogger.Printf("%v \n", v)
}

func Error(v ...interface{}) {
	if !inited {
		logger_init()
	}
	errorLogger.Printf("%v \n", v)
}

func Fatal(v ...interface{}) {
	if !inited {
		logger_init()
	}
	fatalLogger.Printf("%v \n", v)
}
