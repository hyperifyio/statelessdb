// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.
//go:build !debuglog
// +build !debuglog

package logs

//go:format Debugf 1 2
func (l *Logger) Debugf(msg string, args ...interface{}) {
}

//go:format Infof 1 2
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.queue <- LogMessage{InfoLogLevel, msg, args}
}

//go:format Warnf 1 2
func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.queue <- LogMessage{WarnLogLevel, msg, args}
}

//go:format Errorf 1 2
func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.queue <- LogMessage{ErrorLogLevel, msg, args}
}
