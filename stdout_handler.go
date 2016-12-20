package logger

import (
	"io/ioutil"
	"log"
	"os"
)

// Logger provides support to write to log files.
type StdOutHandler struct {
	namespace string

	// Debug is for full detailed messages.
	DebugLogger *log.Logger

	// Info is for important messages.
	InfoLogger *log.Logger

	// Warning is for need to know issue messages.
	WarnLogger *log.Logger

	// Error is for error messages.
	ErrorLogger *log.Logger

	// Fatal
	FatalLogger *log.Logger
}

func NewStdOutHandler(namespace string) *StdOutHandler {
	return &StdOutHandler{
		namespace: namespace,
	}
}

func (handler *StdOutHandler) Init(namespace string, level Level) {
	handler.turnOnLogging(namespace, level)
}

func (handler *StdOutHandler) SetLevel(level Level) {
	handler.turnOnLogging(handler.namespace, level)
}

func (handler *StdOutHandler) turnOnLogging(namespace string, level Level) {
	debugOutput := ioutil.Discard
	infoOutput := ioutil.Discard
	warnOutput := ioutil.Discard
	errorOutput := ioutil.Discard

	switch level {
	case LevelDebug:
		debugOutput, infoOutput, warnOutput = os.Stdout, os.Stdout, os.Stdout
		errorOutput = os.Stderr
	case LevelInfo:
		infoOutput, warnOutput = os.Stdout, os.Stdout
		errorOutput = os.Stderr
	case LevelWarn:
		warnOutput = os.Stdout
		errorOutput = os.Stderr
	case LevelError:
		errorOutput = os.Stderr
	}

	if namespace != "" {
		namespace = "<" + namespace + "> "
	}

	handler.DebugLogger = log.New(debugOutput, namespace+"[DEBUG] ", 0)
	handler.InfoLogger = log.New(infoOutput, namespace+"[INFO] ", 0)
	handler.WarnLogger = log.New(warnOutput, namespace+"[WARN] ", 0)
	handler.ErrorLogger = log.New(errorOutput, namespace+"[ERROR] ", 0)
	handler.FatalLogger = log.New(errorOutput, namespace+"[FATAL] ", 0)
}

func (handler *StdOutHandler) Debug(msg string) {
	handler.DebugLogger.Println(msg)
}

func (handler *StdOutHandler) Info(msg string) {
	handler.InfoLogger.Println(msg)
}

func (handler *StdOutHandler) Warn(msg string) {
	handler.WarnLogger.Println(msg)
}

func (handler *StdOutHandler) Error(msg string) {
	handler.ErrorLogger.Println(msg)
}

func (handler *StdOutHandler) Fatal(msg string) {
	handler.FatalLogger.Println(msg)
}
