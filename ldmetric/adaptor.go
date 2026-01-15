/*
 * Copyright (C) distroy
 */

package ldmetric

import "github.com/distroy/ldgo-base/ldatomic"

var (
	adaptor ldatomic.Any[Adaptor]

	metricPrefix ldatomic.String
)

func SetMetricPrefix(prefix string) { metricPrefix.Store(prefix) }
func getMetricPrefix() string       { return metricPrefix.Load() }

type MetricInfo struct {
	Type       string
	Metric     string
	Action     string              // for `gauge` type of promethues
	Buckets    []float64           // for `histogram` type of promethues
	Objectives map[float64]float64 // for `summary` type of promethues
}

type Adaptor interface {
	GetReportFunc(mi *MetricInfo) (_report func(v float64, labels map[string]string), _reset func())
}

func SetAdaptor(a Adaptor) { adaptor.Store(a) }
func getAdaptor() Adaptor  { return adaptor.Load() }
