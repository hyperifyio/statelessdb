// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package main

import (
	"github.com/hyperifyio/statelessdb/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Declare global metrics
var (
	ComputeDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "compute_duration_milliseconds",
			Help:    "Duration in milliseconds that resources were computed",
			Buckets: prometheus.LinearBuckets(0, 1000, 300),
		},
	)

	ResourceCreatedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "compute_started_total",
			Help: "Count of created resources",
		},
		[]string{}, // Labels
	)
)

func init() {
	metrics.MustRegister(
		ComputeDuration,
		ResourceCreatedTotal,
	)
}

func RecordResourceCreatedMetric() {
	ResourceCreatedTotal.WithLabelValues().Inc()
}
