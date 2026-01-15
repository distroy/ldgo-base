/*
 * Copyright (C) distroy
 */

package ldmetric

import "github.com/distroy/ldgo-base/ldatomic"

type testAdaptor struct {
	metricInfo *MetricInfo
	labels     map[string]string
	value      float64
	reset      ldatomic.Int
}

func (a *testAdaptor) GetReportFunc(mi *MetricInfo) (_report func(v float64, labels map[string]string), _reset func()) {
	a.metricInfo = mi
	_report = func(v float64, labels map[string]string) {
		a.labels = labels
		a.value = v
	}
	_reset = func() {
		a.reset.Add(1)
	}
	return _report, _reset
}
