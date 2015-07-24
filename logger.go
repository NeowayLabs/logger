package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var DefaultLogger *Logger = Namespace("")

type Level uint

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	Namespace   string
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
}

func getEnvVarLevel(namespace string) string {
	prefix := "SEVERINO_LOGGER"
	if namespace != "" {
		prefix += "_"
		namespace = strings.ToUpper(namespace)
		namespace = strings.Replace(namespace, "-", "_", -1)
		namespace = strings.Replace(namespace, ".", "_", -1)
	}

	return strings.ToLower(os.Getenv(prefix + namespace))
}

func getLevelByString(level string) Level {
	if strings.EqualFold(level, "debug") {
		return LevelDebug
	} else if strings.EqualFold(level, "info") {
		return LevelInfo
	} else if strings.EqualFold(level, "warn") {
		return LevelWarn
	} else if strings.EqualFold(level, "error") {
		return LevelError
	} else {
		return LevelInfo
	}
}

func Namespace(namespace string) *Logger {
	logger := &Logger{
		Namespace: namespace,
	}

	level := getEnvVarLevel(namespace)
	logger.SetLevel(getLevelByString(level))

	return logger
}

func (this *Logger) SetLevel(level Level) {
	namespace := this.Namespace
	if namespace != "" {
		namespace = "<" + namespace + "> "
	}

	var debugHandle, infoHandle, warnHandle io.Writer
	if level == LevelDebug {
		debugHandle, infoHandle, warnHandle = os.Stdout, os.Stdout, os.Stdout
	} else if level == LevelInfo {
		debugHandle, infoHandle, warnHandle = ioutil.Discard, os.Stdout, os.Stdout
	} else if level == LevelWarn {
		debugHandle, infoHandle, warnHandle = ioutil.Discard, ioutil.Discard, os.Stdout
	} else if level == LevelError {
		debugHandle, infoHandle, warnHandle = ioutil.Discard, ioutil.Discard, ioutil.Discard
	}

	this.DebugLogger = log.New(debugHandle, namespace+"[DEBUG] ", 0)
	this.InfoLogger = log.New(infoHandle, namespace+"[INFO] ", 0)
	this.WarnLogger = log.New(warnHandle, namespace+"[WARN] ", 0)
	this.ErrorLogger = log.New(os.Stderr, namespace+"[ERROR] ", 0)
}

func (this *Logger) Debug(format string, v ...interface{}) {
	this.DebugLogger.Printf(format+"\n", v...)
}

func (this *Logger) Info(format string, v ...interface{}) {
	this.InfoLogger.Printf(format+"\n", v...)
}

func (this *Logger) Warn(format string, v ...interface{}) {
	this.WarnLogger.Printf(format+"\n", v...)
}

func (this *Logger) Error(format string, v ...interface{}) {
	this.ErrorLogger.Printf(format+"\n", v...)
}

func SetLevel(level Level) {
	DefaultLogger.SetLevel(level)
}

func Debug(format string, v ...interface{}) {
	DefaultLogger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	DefaultLogger.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	DefaultLogger.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	DefaultLogger.Error(format, v...)
}
