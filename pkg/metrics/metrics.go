// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Declare a global counter
var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests",
		},
		[]string{"path"}, // Labels
	)

	ResourceCreatedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "compute_started_total",
			Help: "Count of created resources",
		},
		[]string{}, // Labels
	)

	ComputeDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "compute_duration_milliseconds",
			Help:    "Duration in milliseconds that resources were computed",
			Buckets: prometheus.LinearBuckets(0, 1000, 300),
		},
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

func init() {
	prometheus.MustRegister(
		HttpRequestsTotal,
		FailedOperationsCounter,
		ResourceCreatedTotal,
		ComputeDuration,
		FailedAttemptsHistogram,
	)
}

func RecordFailedOperationMetric(operationName string) {
	// Increment the counter for the specific operation that failed
	FailedOperationsCounter.WithLabelValues(operationName).Inc()
}
