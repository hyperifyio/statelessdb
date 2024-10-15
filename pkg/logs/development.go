// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.
//go:build debuglog
// +build debuglog

package logs

import (
	"log"
)

//go:format Debugf 1 2
func (l *Logger) Debugf(msg string, args ...interface{}) {
	callDepth := getCallStackDepth()
	if !l.visible(DebugLogLevel, callDepth) {
		return
	}
	m := LogMessage{l.context, DebugLogLevel, msg, args, callDepth}
	//l.queue <- m
	log.Println(m.String())
}

//go:format Infof 1 2
func (l *Logger) Infof(msg string, args ...interface{}) {
	callDepth := getCallStackDepth()
	if !l.visible(InfoLogLevel, callDepth) {
		return
	}
	m := LogMessage{l.context, InfoLogLevel, msg, args, callDepth}
	//l.queue <- m
	log.Println(m.String())
}

//go:format Warnf 1 2
func (l *Logger) Warnf(msg string, args ...interface{}) {
	callDepth := getCallStackDepth()
	if !l.visible(WarnLogLevel, callDepth) {
		return
	}
	m := LogMessage{l.context, WarnLogLevel, msg, args, callDepth}
	//l.queue <- m
	log.Println(m.String())
}

//go:format Errorf 1 2
func (l *Logger) Errorf(msg string, args ...interface{}) {
	callDepth := getCallStackDepth()
	if !l.visible(ErrorLogLevel, callDepth) {
		return
	}
	m := LogMessage{l.context, ErrorLogLevel, msg, args, callDepth}
	//l.queue <- m
	log.Println(m.String())
}

// Start starts the log processing goroutine
func (l *Logger) Start() {
}
