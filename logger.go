package logger

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

// DefaultLogger default logger
var DefaultLogger = Namespace("")
var loggers = map[string]*Logger{}
var loggersLock sync.Mutex

// defaultEnvironmentVariablePrefix default environment variable prefix
var defaultEnvironmentVariablePrefixLevel = "LOGGER_LEVEL"
var defaultEnvironmentVariablePrefixOutput = "LOGGER_OUTPUT"

const (
	// LevelNone ...
	LevelNone Level = iota

	// LevelError logs just Errors
	LevelError

	// LevelWarn logs Warning and Errors
	LevelWarn

	// LevelInfo logs Info, Warnings and Errors
	LevelInfo

	// LevelDebug logs everything
	LevelDebug

	// Output stdout
	OutputStdOut Output = iota

	// Output syslog
	OutputSysLog
)

type (
	// Level ...
	Level uint

	// Output ...
	Output uint

	// InitInterface ...
	HandlerInterface interface {
		Init(namespace string, level Level)
		SetLevel(level Level)
	}

	// DebugInterface ...
	DebugInterface interface {
		Debug(msg string)
	}

	// InfoInterface ...
	InfoInterface interface {
		Info(msg string)
	}

	// WarnInterface ...
	WarnInterface interface {
		Warn(msg string)
	}

	// ErrorInterface ...
	ErrorInterface interface {
		Error(msg string)
	}

	// FatalInterface ...
	FatalInterface interface {
		Fatal(msg string)
	}

	// Logger ...
	Logger struct {
		Namespace string
		Handlers  []interface{}
		Level     Level
	}
)

func formatPrefix(namespace string, env string) string {
	if namespace != "" {
		namespace = strings.ToUpper(namespace)
		namespace = strings.Replace(namespace, "-", "_", -1)
		namespace = strings.Replace(namespace, ".", "_", -1)
		return namespace + "_" + env
	}
	return env
}

func getEnvVarLevel(namespace string) string {
	prefix := formatPrefix(namespace, defaultEnvironmentVariablePrefixLevel)

	level := os.Getenv(prefix)
	if level == "" {
		level = os.Getenv(defaultEnvironmentVariablePrefixLevel)
	}

	return strings.ToLower(level)
}

func getEnvVarOutput(namespace string) string {
	prefix := formatPrefix(namespace, defaultEnvironmentVariablePrefixOutput)

	output := os.Getenv(prefix)

	if strings.TrimSpace(output) == "" {
		output = os.Getenv(defaultEnvironmentVariablePrefixOutput)
	}

	return output
}

func setEnvironmentVariablePrefix(prefix string) error {
	loggersLock.Lock()
	defer loggersLock.Unlock()

	for namespace := range loggers {
		if namespace != "" {
			return errors.New("Cannot change prefix because some logs have already been use")
		}
	}

	delete(loggers, "")
	defaultEnvironmentVariablePrefixLevel = prefix

	return nil
}

func SetDefaultEnvironmentVariablePrefix(prefix string) error {
	if err := setEnvironmentVariablePrefix(prefix); err != nil {
		return err
	}
	DefaultLogger = Namespace("")

	return nil
}

func GetDefaultEnvironmentVariablePrefix() string {
	return defaultEnvironmentVariablePrefixLevel
}

// GetLevelByString ...
func GetLevelByString(envLevel string) Level {
	var level Level
	envLevel = strings.ToUpper(envLevel)

	switch envLevel {
	case "DEBUG":
		level = LevelDebug
	case "INFO":
		level = LevelInfo
	case "WARN":
		level = LevelWarn
	case "ERROR":
		level = LevelError
	case "NONE":
		level = LevelNone
	default:
		return LevelInfo
	}

	return level
}

// GetOutputByString ...
func GetOutputByString(envOutput string) Output {
	envOutput = strings.ToUpper(envOutput)

	switch envOutput {
	case "SYSLOG":
		return OutputSysLog
	default:
		return OutputStdOut
	}
}

// Namespace create a new logger namespace (new instance of logger)
func Namespace(namespace string) *Logger {
	loggersLock.Lock()
	defer loggersLock.Unlock()

	namespaceLower := strings.ToLower(namespace)
	if logger, ok := loggers[namespaceLower]; ok {
		return logger
	}

	logger := &Logger{
		Namespace: namespace,
		Level:     GetLevelByString(getEnvVarLevel(namespace)),
	}

	output := GetOutputByString(getEnvVarOutput(namespace))
	if output == OutputSysLog {
		logger.AddHandler(NewSysLogHandler(namespace))
	} else {
		logger.AddHandler(NewStdOutHandler(namespace))
	}

	loggers[namespaceLower] = logger

	return logger
}

// AddHandler ...
func (logger *Logger) AddHandler(handler interface{}) {
	logger.Handlers = append(logger.Handlers, handler)

	if handlerInterface, ok := handler.(HandlerInterface); ok {
		handlerInterface.Init(logger.Namespace, logger.Level)
	}
}

// SetLevel ...
func (logger *Logger) SetLevel(level Level) {
	logger.Level = level

	for _, handler := range logger.Handlers {
		if handlerInterface, ok := handler.(HandlerInterface); ok {
			handlerInterface.SetLevel(level)
		}
	}
}

// Debug ...
func (logger *Logger) Debug(format string, v ...interface{}) {
	if logger.Level < LevelDebug {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if debugHandler, ok := handler.(DebugInterface); ok {
			debugHandler.Debug(msg)
		}
	}
}

// Info ...
func (logger *Logger) Info(format string, v ...interface{}) {
	if logger.Level < LevelInfo {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if infoHandler, ok := handler.(InfoInterface); ok {
			infoHandler.Info(msg)
		}
	}
}

// Warn ...
func (logger *Logger) Warn(format string, v ...interface{}) {
	if logger.Level < LevelWarn {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if warnHandler, ok := handler.(WarnInterface); ok {
			warnHandler.Warn(msg)
		}
	}
}

// Error ...
func (logger *Logger) Error(format string, v ...interface{}) {
	if logger.Level < LevelError {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if errorHandler, ok := handler.(ErrorInterface); ok {
			errorHandler.Error(msg)
		}
	}
}

// Fatal ...
func (logger *Logger) Fatal(format string, v ...interface{}) {
	if logger.Level < LevelError {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if fatalHandler, ok := handler.(FatalInterface); ok {
			fatalHandler.Fatal(msg)
		}
	}
	os.Exit(1)
}

// Write ...
func (logger *Logger) Write(b []byte) (int, error) {
	logger.Info("%s", strings.TrimRight(string(b), "\n"))
	return len(b), nil
}

// AddHandler ...
func AddHandler(handler interface{}) {
	DefaultLogger.AddHandler(handler)
}

// SetLevel ...
func SetLevel(level Level) {
	DefaultLogger.SetLevel(level)
}

// Debug ...
func Debug(format string, v ...interface{}) {
	DefaultLogger.Debug(format, v...)
}

// Info ...
func Info(format string, v ...interface{}) {
	DefaultLogger.Info(format, v...)
}

// Warn ...
func Warn(format string, v ...interface{}) {
	DefaultLogger.Warn(format, v...)
}

// Error ...
func Error(format string, v ...interface{}) {
	DefaultLogger.Error(format, v...)
}

// Fatal ...
func Fatal(format string, v ...interface{}) {
	DefaultLogger.Fatal(format, v...)
}
