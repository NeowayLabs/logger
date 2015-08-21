package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

type (
	DefaultHandler struct {
		DebugLogger *log.Logger
		InfoLogger  *log.Logger
		WarnLogger  *log.Logger
		ErrorLogger *log.Logger
		FatalLogger *log.Logger
	}
)

func (handler *DefaultHandler) Init(namespace string, level Level) {
	if namespace != "" {
		namespace = "<" + namespace + "> "
	}

	var debugOutput, infoOutput, warnOutput io.Writer
	if level == LevelDebug {
		debugOutput, infoOutput, warnOutput = os.Stdout, os.Stdout, os.Stdout
	} else if level == LevelInfo {
		debugOutput = ioutil.Discard
		infoOutput, warnOutput = os.Stdout, os.Stdout
	} else {
		debugOutput, infoOutput = ioutil.Discard, ioutil.Discard
		warnOutput = os.Stdout
	}

	handler.DebugLogger = log.New(debugOutput, namespace+"[DEBUG] ", 0)
	handler.InfoLogger = log.New(infoOutput, namespace+"[INFO] ", 0)
	handler.WarnLogger = log.New(warnOutput, namespace+"[WARN] ", 0)
	handler.ErrorLogger = log.New(os.Stderr, namespace+"[ERROR] ", 0)
	handler.FatalLogger = log.New(os.Stderr, namespace+"[FATAL] ", 0)
}

func (handler *DefaultHandler) Debug(msg string) {
	handler.DebugLogger.Println(msg)
}

func (handler *DefaultHandler) Info(msg string) {
	handler.InfoLogger.Println(msg)
}

func (handler *DefaultHandler) Warn(msg string) {
	handler.WarnLogger.Println(msg)
}

func (handler *DefaultHandler) Error(msg string) {
	handler.ErrorLogger.Println(msg)
}

func (handler *DefaultHandler) Fatal(msg string) {
	handler.FatalLogger.Println(msg)
}
