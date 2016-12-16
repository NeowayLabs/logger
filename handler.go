package logger

import (
	"io"
	"io/ioutil"
	"log"
	"log/syslog"
	"os"
	"strings"
)

// Logger provides support to write to log files.
type DefaultHandler struct {
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

func (handler *DefaultHandler) Init(namespace string, level Level) {
	handler.turnOnLogging(namespace, level, nil)
}

func (handler *DefaultHandler) InitSysLog(namespace, address string, level Level) {
	var sysLogLevel syslog.Priority
	network := "udp"

	switch level {
	case LevelDebug:
		sysLogLevel = syslog.LOG_DEBUG
	case LevelInfo:
		sysLogLevel = syslog.LOG_INFO
	case LevelWarn:
		sysLogLevel = syslog.LOG_WARNING
	case LevelError:
		sysLogLevel = syslog.LOG_ERR
	}

	addressSplit := strings.Split(address, "/")
	address = addressSplit[0]
	if len(addressSplit) == 2 {
		network = addressSplit[1]
	}

	sysLogHandle, err := syslog.Dial(network, address, sysLogLevel, namespace)
	if err != nil {
		log.Fatal("error")
	}

	handler.turnOnLogging(namespace, level, sysLogHandle)
}

func (handler *DefaultHandler) turnOnLogging(namespace string, level Level, sysLogHandle io.Writer) {
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

	if sysLogHandle != nil {
		if debugOutput == os.Stdout {
			debugOutput = sysLogHandle
		}

		if infoOutput == os.Stdout {
			infoOutput = sysLogHandle
		}

		if warnOutput == os.Stdout {
			warnOutput = sysLogHandle
		}

		if errorOutput == os.Stderr {
			errorOutput = sysLogHandle
		}
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
