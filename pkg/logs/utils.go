// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package logs

import (
	"runtime"
	"strings"
)

func getCallStackDepth() int {
	// Capture the call stack up to 32 frames deep
	var pcs [32]uintptr
	n := runtime.Callers(2, pcs[:])
	return n
}

func trimFilePath(path string) string {
	// Trims the path to get the last portion for readability
	slash := strings.LastIndex(path, "/")
	if slash == -1 {
		return path
	}
	return path[slash+1:]
}
