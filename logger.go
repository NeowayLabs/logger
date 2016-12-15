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
var defaultEnvironmentVariablePrefix = "SEVERINO_LOGGER"
var defaultEnvironmentVariablePrefixSysLog = "SEVERINO_SYSLOG"
var defaultEnvironmentVariablePrefixSysLogNetwork = "SEVERINO_SYSLOG_NETWORK"
var defaultEnvironmentVariablePrefixSysLogAddress = "SEVERINO_SYSLOG_ADDRESS"

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
)

type (
	// Level ...
	Level uint
	// Interface ...
	Interface interface {
	}
	// InitInterface ...
	InitInterface interface {
		Init(namespace string, level Level)
		InitSysLog(namespace, network, address string, level Level)
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
		Namespace     string
		Level         Level
		Handlers      []Interface
		IsSysLog      bool
		SysLogAddress string
		SysLogNetwork string
	}
)

func getEnvVarLevel(namespace string) string {
	prefix := defaultEnvironmentVariablePrefix
	if namespace != "" {
		prefix += "_"
		namespace = strings.ToUpper(namespace)
		namespace = strings.Replace(namespace, "-", "_", -1)
		namespace = strings.Replace(namespace, ".", "_", -1)
	}

	level := os.Getenv(prefix + namespace)
	if level == "" {
		level = os.Getenv(defaultEnvironmentVariablePrefix)
	}

	return strings.ToLower(level)
}

func getEnvVarSysLog(namespace string) string {
	prefix := defaultEnvironmentVariablePrefixSysLog
	if namespace != "" {
		prefix += "_"
		namespace = strings.ToUpper(namespace)
		namespace = strings.Replace(namespace, "-", "_", -1)
		namespace = strings.Replace(namespace, ".", "_", -1)
	}

	sysLog := os.Getenv(prefix + namespace)
	if sysLog == "" {
		sysLog = os.Getenv(defaultEnvironmentVariablePrefix)
	}

	return sysLog
}

func getEnvVarSysLogAddress(namespace string) string {
	prefix := defaultEnvironmentVariablePrefixSysLogAddress
	if namespace != "" {
		prefix += "_"
		namespace = strings.ToUpper(namespace)
		namespace = strings.Replace(namespace, "-", "_", -1)
		namespace = strings.Replace(namespace, ".", "_", -1)
	}

	address := os.Getenv(prefix + namespace)
	if address == "" {
		address = os.Getenv(defaultEnvironmentVariablePrefixSysLogAddress)
	}

	return address
}

func getEnvVarSysLogNetwork(namespace string) string {
	prefix := defaultEnvironmentVariablePrefixSysLogNetwork
	if namespace != "" {
		prefix += "_"
		namespace = strings.ToUpper(namespace)
		namespace = strings.Replace(namespace, "-", "_", -1)
		namespace = strings.Replace(namespace, ".", "_", -1)
	}

	network := os.Getenv(prefix + namespace)
	if network == "" {
		network = os.Getenv(defaultEnvironmentVariablePrefixSysLogNetwork)
	}

	return network
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
	defaultEnvironmentVariablePrefix = prefix

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
	return defaultEnvironmentVariablePrefix
}

// GetLevelByString ...
func GetLevelByString(envLevel string) Level {
	var level Level
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
		level = LevelError
	default:
		return LevelInfo
	}
	return level
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
		IsSysLog:  getEnvVarSysLog(namespace) == true,
	}

	logger.SetLevel(GetLevelByString(getEnvVarLevel(namespace)))
	logger.AddHandler(&DefaultHandler{})

	loggers[namespaceLower] = logger

	if logger.IsSysLog {
		logger.SysLogNetwork = getEnvVarSysLogNetwork(namespace)
		logger.SysLogAddress = getEnvVarSysLogAddress(namespace)
	}

	return logger
}

// AddHandler ...
func (logger *Logger) AddHandler(handler Interface) {
	logger.Handlers = append(logger.Handlers, handler)

	if initHandler, ok := handler.(InitInterface); ok {
		initHandler.Init(logger.Namespace, logger.Level)
	}
}

// SetLevel ...
func (logger *Logger) SetLevel(level Level) {
	logger.Level = level

	for _, handler := range logger.Handlers {
		if initHandler, ok := handler.(InitInterface); ok {
			if logger.IsSysLog {
				initHandler.InitSysLog(logger.Namespace, logger.SysLogNetwork, logger.SysLogAddress, logger.Level)
				continue
			}
			initHandler.Init(logger.Namespace, logger.Level)
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
func AddHandler(handler Interface) {
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
