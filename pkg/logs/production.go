// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.
//go:build !debuglog
// +build !debuglog

package logs

import "log"

//go:format Debugf 1 2
func (l *Logger) Debugf(msg string, args ...interface{}) {
}

//go:format Infof 1 2
func (l *Logger) Infof(msg string, args ...interface{}) {
	callDepth := getCallStackDepth()
	l.queue <- LogMessage{l.context, InfoLogLevel, msg, args, callDepth}
}

//go:format Warnf 1 2
func (l *Logger) Warnf(msg string, args ...interface{}) {
	callDepth := getCallStackDepth()
	l.queue <- LogMessage{l.context, WarnLogLevel, msg, args, callDepth}
}

//go:format Errorf 1 2
func (l *Logger) Errorf(msg string, args ...interface{}) {
	callDepth := getCallStackDepth()
	l.queue <- LogMessage{l.context, ErrorLogLevel, msg, args, callDepth}
}

// Start starts the log processing goroutine
func (l *Logger) Start() {
	l.wg.Add(1)
	go l.processQueue()
}

func (l *Logger) Stop() {
	l.closeOnce.Do(func() {
		close(l.queue)
		l.wg.Wait()
	})
}

func (l *Logger) processQueue() {
	defer l.wg.Done()
	for msg := range l.queue {
		if l.visible(msg.level, msg.depth) {
			log.Println(msg.String())
		}
	}
}
