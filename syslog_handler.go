package logger

import (
	"io"
	"io/ioutil"
	"log"
	"log/syslog"
	"os"
	"strings"
)

var defaultEnvironmentVariablePrefixSysLogAddress = "LOGGER_SYSLOG_ADDRESS"

// Logger provides support to write to log files.
type SysLogHandler struct {
	namespace    string
	network      string
	address      string
	sysLogWriter *syslog.Writer

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

func NewSysLogHandler(namespace string) *SysLogHandler {
	network := "udp"

	syslogAddress := getEnvVarSysLogAddress(namespace)
	if strings.TrimSpace(syslogAddress) == "" {
		log.Fatal("To use SysLog, the environment variable LOGGER_SYSLOG_ADDRESS is required.")
	}
	addressSplit := strings.Split(syslogAddress, "/")
	address := addressSplit[0]
	if len(addressSplit) == 2 {
		network = addressSplit[1]
	}

	return &SysLogHandler{
		namespace: namespace,
		network:   network,
		address:   address,
	}
}

func (handler *SysLogHandler) Init(namespace string, level Level) {
	handler.sysLogWriter = dial(handler.namespace, handler.network, handler.address, getPriority(level))

	handler.turnOnLogging(namespace, level, handler.sysLogWriter)
}

func (handler *SysLogHandler) SetLevel(level Level) {
	if handler.sysLogWriter != nil {
		handler.sysLogWriter.Close()
	}

	handler.sysLogWriter = dial(handler.namespace, handler.network, handler.address, getPriority(level))

	handler.turnOnLogging(handler.namespace, level, handler.sysLogWriter)
}

func dial(namespace, network, address string, priority syslog.Priority) *syslog.Writer {
	sysLogWriter, err := syslog.Dial(network, address, priority, namespace)
	if err != nil {
		log.Fatalf("Error on syslog.Dial(%s, %s, %s, %s)\n", network, address, priority, namespace)
	}
	return sysLogWriter
}

func getPriority(level Level) syslog.Priority {
	var priority syslog.Priority
	switch level {
	case LevelDebug:
		priority = syslog.LOG_DEBUG
	case LevelInfo:
		priority = syslog.LOG_INFO
	case LevelWarn:
		priority = syslog.LOG_WARNING
	case LevelError:
		priority = syslog.LOG_ERR
	}
	return priority
}

func getEnvVarSysLogAddress(namespace string) string {
	prefix := defaultEnvironmentVariablePrefixSysLogAddress

	if namespace != "" {
		namespace = strings.ToUpper(namespace)
		namespace = strings.Replace(namespace, "-", "_", -1)
		namespace = strings.Replace(namespace, ".", "_", -1)
		prefix = namespace + "_" + prefix

	}

	address := os.Getenv(prefix)
	if address == "" {
		address = os.Getenv(defaultEnvironmentVariablePrefixSysLogAddress)
	}

	return address
}

func (handler *SysLogHandler) turnOnLogging(namespace string, level Level, sysLogWriter io.Writer) {
	debugOutput := ioutil.Discard
	infoOutput := ioutil.Discard
	warnOutput := ioutil.Discard
	errorOutput := ioutil.Discard

	switch level {
	case LevelDebug:
		debugOutput, infoOutput, warnOutput, errorOutput = sysLogWriter, sysLogWriter, sysLogWriter, sysLogWriter
	case LevelInfo:
		infoOutput, warnOutput, errorOutput = sysLogWriter, sysLogWriter, sysLogWriter
	case LevelWarn:
		warnOutput, errorOutput = sysLogWriter, sysLogWriter
	case LevelError:
		errorOutput = sysLogWriter
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

func (handler *SysLogHandler) Debug(msg string) {
	handler.DebugLogger.Println(msg)
}

func (handler *SysLogHandler) Info(msg string) {
	handler.InfoLogger.Println(msg)
}

func (handler *SysLogHandler) Warn(msg string) {
	handler.WarnLogger.Println(msg)
}

func (handler *SysLogHandler) Error(msg string) {
	handler.ErrorLogger.Println(msg)
}

func (handler *SysLogHandler) Fatal(msg string) {
	handler.FatalLogger.Println(msg)
}
