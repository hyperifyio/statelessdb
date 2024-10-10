// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package apis

import (
	"net/http"
	"statelessdb/pkg/metrics"
)

// Make sure these are used only once per place, e.g. IDE should report 1 usage
// for each! Create another constant for each error.

const (
	EncodingFailedError    = "encoding-failed"
	WritingBodyFailedError = "writing-body-failed"
	EncryptionFailedError  = "encryption-failed"
	DecryptionFailedError  = "decryption-failed"
	BadPrivateBodyError    = "bad-private-body"
	BadBodyError           = "bad-body"
	ComputeLogicError      = "compute-logic-error"
)

func sendHttpError(w http.ResponseWriter, code string, status int) {
	metrics.RecordFailedOperationMetric(code)
	http.Error(w, code, status)
}
