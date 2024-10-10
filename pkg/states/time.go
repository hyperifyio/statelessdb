// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package states

import "time"

func NewTimeNow() int64 {
	return time.Now().UnixMilli()
}
