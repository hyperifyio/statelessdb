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
	if callDepth > l.depth {
		return
	}

	m := LogMessage{DebugLogLevel, msg, args, callDepth}
	//l.queue <- m
	log.Println("[" + l.Context + "] " + m.String())
}

//go:format Infof 1 2
func (l *Logger) Infof(msg string, args ...interface{}) {
	callDepth := getCallStackDepth()
	if callDepth > l.depth {
		return
	}
	m := LogMessage{InfoLogLevel, msg, args, callDepth}
	//l.queue <- m
	log.Println("[" + l.Context + "] " + m.String())
}

//go:format Warnf 1 2
func (l *Logger) Warnf(msg string, args ...interface{}) {
	callDepth := getCallStackDepth()
	if callDepth > l.depth {
		return
	}
	m := LogMessage{WarnLogLevel, msg, args, callDepth}
	//l.queue <- m
	log.Println("[" + l.Context + "] " + m.String())
}

//go:format Errorf 1 2
func (l *Logger) Errorf(msg string, args ...interface{}) {
	callDepth := getCallStackDepth()
	if callDepth > l.depth {
		return
	}
	m := LogMessage{ErrorLogLevel, msg, args, callDepth}
	//l.queue <- m
	log.Println("[" + l.Context + "] " + m.String())
}
