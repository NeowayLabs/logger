package logger

import (
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

	handler.DebugLogger = log.New(os.Stdout, namespace+"[DEBUG] ", 0)
	handler.InfoLogger = log.New(os.Stdout, namespace+"[INFO] ", 0)
	handler.WarnLogger = log.New(os.Stdout, namespace+"[WARN] ", 0)
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
