/*
 * Copyright (C) distroy
 */

package ldmetric

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

type testMetricWithMethod struct {
	Index int64
	Queue string
}

func (_ testMetricWithMethod) MetricType() string       { return "meter" }
func (_ testMetricWithMethod) MetricName() string       { return "meter_normal_metric" }
func (_ testMetricWithMethod) MetricAction() string     { return "store" }
func (_ testMetricWithMethod) MetricBuckets() []float64 { return []float64{1} }
func (_ testMetricWithMethod) MetricObjectives() map[float64]float64 {
	return map[float64]float64{1: 1}
}

func TestReport(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		adaptor := &testAdaptor{}
		SetAdaptor(adaptor)
		SetMetricPrefix("test.")

		c.So(adaptor, convey.ShouldResemble, &testAdaptor{
			metricInfo: nil,
			labels:     nil,
			value:      0,
			reset:      0,
		})

		c.Convey("with method", func(c convey.C) {
			Report(&testMetricWithMethod{
				Index: 3,
				Queue: "test_queue",
			})
			c.So(adaptor, convey.ShouldResemble, &testAdaptor{
				metricInfo: &MetricInfo{
					Type:       "meter",
					Metric:     "meter_normal_metric",
					Action:     "store",
					Buckets:    []float64{1},
					Objectives: map[float64]float64{1: 1},
				},
				labels: map[string]string{
					`index`: "3",
					`queue`: "test_queue",
				},
				value: 1,
				reset: 0,
			})
			ResetReporter(&testMetricWithMethod{})
			c.So(adaptor, convey.ShouldResemble, &testAdaptor{
				metricInfo: &MetricInfo{
					Type:       "meter",
					Metric:     "meter_normal_metric",
					Action:     "store",
					Buckets:    []float64{1},
					Objectives: map[float64]float64{1: 1},
				},
				labels: map[string]string{
					`index`: "3",
					`queue`: "test_queue",
				},
				value: 1,
				reset: 1,
			})
		})

		c.Convey("with field tag", func(c convey.C) {
			type testMetric struct {
				Count float64 `ldmetric:"type:counter; metric:nomarl_metric; action:set"`
				Index int64
			}

			Report(&testMetric{
				Count: 2,
				Index: 3,
			})
			c.So(adaptor, convey.ShouldResemble, &testAdaptor{
				metricInfo: &MetricInfo{
					Type:   "counter",
					Metric: "nomarl_metric",
					Action: "set",
				},
				labels: map[string]string{
					`index`: "3",
				},
				value: 2,
				reset: 0,
			})
			ResetReporter(&testMetric{})
			c.So(adaptor, convey.ShouldResemble, &testAdaptor{
				metricInfo: &MetricInfo{
					Type:   "counter",
					Metric: "nomarl_metric",
					Action: "set",
				},
				labels: map[string]string{
					`index`: "3",
				},
				value: 2,
				reset: 1,
			})
		})

		c.Convey("no field tag & no method", func(c convey.C) {
			type NoFieldTagMetric struct {
				Index int64
			}

			Report(&NoFieldTagMetric{
				Index: 3,
			})
			c.So(adaptor, convey.ShouldResemble, &testAdaptor{
				metricInfo: &MetricInfo{
					Type:   "",
					Metric: "test.no_field_tag_metric",
				},
				labels: map[string]string{
					`index`: "3",
				},
				value: 1,
				reset: 0,
			})
			ResetReporter(&NoFieldTagMetric{})
			c.So(adaptor, convey.ShouldResemble, &testAdaptor{
				metricInfo: &MetricInfo{
					Type:   "",
					Metric: "test.no_field_tag_metric",
				},
				labels: map[string]string{
					`index`: "3",
				},
				value: 1,
				reset: 1,
			})
		})
	})
}
