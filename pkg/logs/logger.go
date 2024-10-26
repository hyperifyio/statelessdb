// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package logs

import (
	"fmt"
	"sync"
)

var (
	DefaultLogLevel   LogLevel = DebugLogLevel
	DefaultBufferSize int      = 2000
	DefaultLogDepth   int      = 10
)

type Logger struct {
	context   string
	level     LogLevel
	queue     chan LogMessage
	wg        sync.WaitGroup
	closeOnce sync.Once
	depth     int
}

func NewLogger(context string) *Logger {
	logger := &Logger{
		context: context,
		level:   DefaultLogLevel,
		queue:   make(chan LogMessage, DefaultBufferSize),
		depth:   DefaultLogDepth,
	}
	logger.Start()
	return logger
}

func (l *Logger) Close() {
	close(l.queue)
}

func (l *Logger) Level(level LogLevel) *Logger {
	l.level = level
	return l
}

func (l *Logger) Depth(depth int) *Logger {
	l.depth = depth
	return l
}

func (l *Logger) String() string {
	return fmt.Sprintf(
		"Logger('%s', '%s', %d)",
		l.context, l.level.String(), l.depth,
	)
}

func (l *Logger) visible(level LogLevel, depth int) bool {
	return level <= l.level && depth <= l.depth
}
