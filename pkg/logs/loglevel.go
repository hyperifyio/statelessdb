// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package logs

type LogLevel int

const (
	NoneLogLevel LogLevel = iota
	ErrorLogLevel
	WarnLogLevel
	InfoLogLevel
	DebugLogLevel
	AllLogLevel
)

func (l LogLevel) String() string {
	switch l {
	case AllLogLevel:
		return "all"
	case DebugLogLevel:
		return "debug"
	case InfoLogLevel:
		return "info"
	case WarnLogLevel:
		return "warn"
	case ErrorLogLevel:
		return "error"
	case NoneLogLevel:
		return "none"
	}
	return "invalid"
}
