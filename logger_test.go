package logger_test

import (
	"sync"
	"testing"

	"github.com/NeowayLabs/logger"
)

func TestCanCreateSameNamespaceOnDifferentGoRoutines(t *testing.T) {
	const concurrency = 10000

	wait := sync.WaitGroup{}
	wait.Add(concurrency)
	createLogger := func() {
		logger.Namespace("test")
		wait.Done()
	}

	for i := 0; i < concurrency; i++ {
		go createLogger()
	}
	wait.Wait()
}
