package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func levelToString(level Level) string {
	if level == LevelDebug {
		return "debug"
	} else if level == LevelInfo {
		return "info"
	} else if level == LevelWarn {
		return "warn"
	} else if level == LevelError {
		return "error"
	} else {
		return ""
	}
}

// HTTPHandler it's a handler to HTTPFunc function
func HTTPHandler() http.Handler {
	return http.HandlerFunc(HTTPFunc)
}

// HTTPFunc permit you control level of all your namespace, and change it in execution time
func HTTPFunc(w http.ResponseWriter, r *http.Request) {
	lastpart := strings.LastIndex(r.RequestURI, "/")
	namespace := r.RequestURI[lastpart+1:]

	// Get list of namespaces and levels
	if r.Method == "GET" {
		// Get all namespaces
		if lastpart == 0 {
			namespaces := make(map[string]string, 0)
			for namespace, logger := range loggers {
				namespace = logger.Namespace
				if namespace == "" {
					namespace = "_default_"
				}
				namespaces[namespace] = levelToString(logger.Level)
			}

			json, _ := json.Marshal(&namespaces)

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			io.WriteString(w, string(json))

			return
		}

		if logger, ok := loggers[namespace]; ok {
			loggerObj := make(map[string]string, 0)
			loggerObj["namespace"] = logger.Namespace
			loggerObj["level"] = levelToString(logger.Level)

			json, _ := json.Marshal(&loggerObj)

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			io.WriteString(w, string(json))
		} else {
			http.Error(w, fmt.Sprintf("namespace '%s' not found", namespace), http.StatusNotFound)
		}

		return
	}

	if r.Method == "PUT" {
		var userLevel map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&userLevel); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if userLevel["level"] == nil {
			http.Error(w, "missing 'level' field", http.StatusBadRequest)
			return
		}

		if lastpart == 0 {
			if userLevel["namespace"] == nil {
				http.Error(w, "missing 'namespace' field", http.StatusBadRequest)
				return
			}

			namespace = userLevel["namespace"].(string)
		}

		level := GetLevelByString(userLevel["level"].(string))
		if logger, ok := loggers[namespace]; ok {
			logger.SetLevel(level)
		} else if namespace == "" {
			DefaultLogger.SetLevel(level)
		} else {
			http.Error(w, fmt.Sprintf("namespace '%s' not found", namespace), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, http.StatusText(http.StatusOK))

		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	io.WriteString(w, http.StatusText(http.StatusNotImplemented))
}
