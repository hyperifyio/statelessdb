// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Declare global metrics
var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests",
		},
		[]string{"path"}, // Labels
	)

	FailedOperationsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "compute_failed_operations_total",
			Help: "Total number of failed operations",
		},
		[]string{"operation"},
	)

	FailedAttemptsHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "compute_failed_attempts",
		Help:    "Histogram of failed attempts",
		Buckets: prometheus.LinearBuckets(0, 10, 50),
	})
)

var (
	ourCollectors = []prometheus.Collector{
		HttpRequestsTotal,
		FailedOperationsCounter,
		FailedAttemptsHistogram,
	}
)

func MustRegister(cs ...prometheus.Collector) {
	our := make([]prometheus.Collector, 0, len(ourCollectors)+len(cs))
	our = append(our, ourCollectors...)
	our = append(our, cs...)
	prometheus.MustRegister(our...)
}

func RecordFailedOperationMetric(operationName string) {
	// Increment the counter for the specific operation that failed
	FailedOperationsCounter.WithLabelValues(operationName).Inc()
}

func RecordHttpRequestMetric(path string) {
	HttpRequestsTotal.WithLabelValues(path).Inc()
}
