// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package logs

import (
	"log"
	"sync"
)

var (
	DefaultLogLevel   LogLevel = DebugLogLevel
	DefaultBufferSize int      = 2000
	DefaultLogDepth   int      = 3
)

type Logger struct {
	Context   string
	Level     LogLevel
	queue     chan LogMessage
	wg        sync.WaitGroup
	closeOnce sync.Once
	depth     int
}

func NewLogger(context string) *Logger {
	logger := &Logger{
		Context: context,
		Level:   DefaultLogLevel,
		queue:   make(chan LogMessage, DefaultBufferSize),
		depth:   DefaultLogDepth,
	}
	// Start the log processing goroutine
	logger.wg.Add(1)
	go logger.processQueue()

	return logger
}

func (l *Logger) String() string {
	return "Logger(" + l.Level.String() + ")"
}

func (l *Logger) processQueue() {
	defer l.wg.Done()
	for msg := range l.queue {
		log.Println("[" + l.Context + "] " + msg.String())
	}
}
