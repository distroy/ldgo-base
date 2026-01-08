/*
 * Copyright (C) distroy
 */

package ldmetric

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/distroy/ldgo-base/lderr"
	"github.com/distroy/ldgo-base/ldstr"
	"github.com/distroy/ldgo-base/ldsync"
	"github.com/distroy/ldgo-base/ldtags"
)

const (
	tagName = "ldmetric"
	nilStr  = ""
)

type stringer interface {
	fmt.Stringer
}

var (
	typeOfError    = reflect.TypeFor[error]()
	typeOfStringer = reflect.TypeFor[stringer]()
	typeOfTime     = reflect.TypeFor[time.Time]()
)

var (
	typePool = &ldsync.Map[reflect.Type, *typeInfo]{}
)

type fieldValueFunc = func(c *toLabelsContext, v reflect.Value)

type typeFieldInfo struct {
	Field  reflect.StructField
	Index  int
	Name   string
	Value  fieldValueFunc
	Layout string
}

type typeMetricField struct {
	Field            reflect.StructField
	Index            int
	Value            func(v reflect.Value) float64
	MetricAction     string
	MetricType       string
	MetricName       string
	MetricObjectives map[float64]float64
	MetricBuckets    []float64
}

type typeInfo struct {
	Type   reflect.Type
	Fields []*typeFieldInfo

	metricField *typeMetricField
	funcInited  ldsync.Once
	reportFunc  func(v reflect.Value, t *typeInfo)
	resetFunc   func()
}

func (t *typeInfo) Report(v reflect.Value) {
	t.funcInited.Do(func() error {
		t.reportFunc, t.resetFunc = t.getReportFunc()
		return nil
	})
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t.reportFunc(v, t)
}

func (t *typeInfo) Reset() {
	if t.funcInited.Done() {
		t.resetFunc()
	}
}

func (t *typeInfo) getReportFunc() (_report func(v reflect.Value, t *typeInfo), _reset func()) {
	mf := t.getMetricInfo()
	if mf.Index < 0 || mf.Value == nil {
		mf.Value = func(_ reflect.Value) float64 { return 1 }
	}

	mi := &MetricInfo{
		Type:       mf.MetricType,
		Metric:     mf.MetricName,
		Action:     mf.MetricAction,
		Buckets:    mf.MetricBuckets,
		Objectives: mf.MetricObjectives,
	}

	adaptor := getAdaptor()
	report, reset := adaptor.GetReportFunc(mi)
	return func(v reflect.Value, t *typeInfo) { report(mf.Value(v), toLabels(v, t)) }, reset

	// switch mf.MetricType {
	// default:
	// 	fallthrough
	// case "counter":
	// 	m := NewCounter(mf.MetricName)
	// 	if mf.Index < 0 || mf.Value == nil {
	// 		return func(v reflect.Value, t *typeInfo) { getClient().EmitCounter(mf.MetricName, 1, toLabels(v, t)) }, func() { m.Reset() }
	// 	} else {
	// 		return func(v reflect.Value, t *typeInfo) { m.With(toLabels(v, t)).Add(mf.Value(v)) }, func() { m.Reset() }
	// 	}
	//
	// case "summary":
	// 	m := NewSummary(mf.MetricName, mf.MetricObjectives)
	// 	return func(v reflect.Value, t *typeInfo) { m.With(toLabels(v, t)).Observe(mf.Value(v)) }, func() { m.Reset() }
	//
	// case "histogram":
	// 	m := NewHistogram(mf.MetricName, mf.MetricBuckets)
	// 	return func(v reflect.Value, t *typeInfo) { m.With(toLabels(v, t)).Observe(mf.Value(v)) }, func() { m.Reset() }
	//
	// case "gauge":
	// 	m := NewGauge(mf.MetricName)
	//
	// 	switch mf.MetricAction {
	// 	case "set":
	// 		return func(v reflect.Value, t *typeInfo) { m.With(toLabels(v, t)).Set(mf.Value(v)) }, func() { m.Reset() }
	//
	// 	case "add":
	// 		return func(v reflect.Value, t *typeInfo) { m.With(toLabels(v, t)).Add(mf.Value(v)) }, func() { m.Reset() }
	//
	// 	case "sub":
	// 		return func(v reflect.Value, t *typeInfo) { m.With(toLabels(v, t)).Sub(mf.Value(v)) }, func() { m.Reset() }
	//
	// 	case "inc":
	// 		return func(v reflect.Value, t *typeInfo) { m.With(toLabels(v, t)).Inc() }, func() { m.Reset() }
	//
	// 	case "dec":
	// 		return func(v reflect.Value, t *typeInfo) { m.With(toLabels(v, t)).Dec() }, func() { m.Reset() }
	// 	}
	//
	// 	if mf.Value == nil {
	// 		return func(v reflect.Value, t *typeInfo) { m.With(toLabels(v, t)).Inc() }, func() { m.Reset() }
	// 	}
	// 	return func(v reflect.Value, t *typeInfo) { m.With(toLabels(v, t)).Set(mf.Value(v)) }, func() { m.Reset() }
	// }
}

func (t *typeInfo) getMetricInfo() *typeMetricField {
	mf := t.metricField
	if mf == nil {
		mf = &typeMetricField{Index: -1}
		t.metricField = mf
	}

	v := reflect.New(t.Type).Interface()

	if vv, _ := v.(interface{ MetricType() string }); vv != nil {
		mf.MetricType = vv.MetricType()
	}

	if vv, _ := v.(interface{ MetricName() string }); vv != nil {
		mf.MetricName = vv.MetricName()
	} else if mf.MetricName == "" {
		mf.MetricName = getMetricPrefix() + ldstr.ToSnakeCase(t.Type.Name())
	}

	if vv, _ := v.(interface{ MetricAction() string }); vv != nil {
		mf.MetricAction = vv.MetricAction()
	}

	if vv, _ := v.(interface{ MetricBuckets() []float64 }); vv != nil {
		mf.MetricBuckets = vv.MetricBuckets()
	}
	if vv, _ := v.(interface{ MetricObjectives() map[float64]float64 }); vv != nil {
		mf.MetricObjectives = vv.MetricObjectives()
	}

	return mf
}

type toLabelsContext struct {
	Name    []string
	Labels  map[string]string
	Current *typeFieldInfo
}

func (c *toLabelsContext) GetName() string { return strings.Join(c.Name, ".") }
func (c *toLabelsContext) Set(v string)    { c.Labels[c.GetName()] = v }

func ToLabels(obj any) map[string]string {
	v := reflect.ValueOf(obj)
	t := getTypeInfo(v.Type())
	return toLabels(v, t)
}

func toLabels(v reflect.Value, t *typeInfo) map[string]string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	c := &toLabelsContext{}
	structToLabels(c, v, t)
	return c.Labels
}

func getTypeInfo(typ reflect.Type) *typeInfo {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("the reporter type should be struct or pointer to struct. %s", typ.String()))
	}

	pool := typePool
	p := pool.Get(typ)
	if p != nil {
		return p
	}

	p = newTypeInfo(typ)
	p, _ = pool.LoadOrStore(typ, p)
	return p
}

func newTypeInfo(typ reflect.Type) *typeInfo {
	n := typ.NumField()
	res := &typeInfo{
		Type:   typ,
		Fields: make([]*typeFieldInfo, 0, n),
	}
	for i := range n {
		sf := typ.Field(i)
		fi, fm := newTypeFieldInfo(i, sf)
		switch {
		case fi != nil:
			res.Fields = append(res.Fields, fi)

		case fm != nil && res.metricField == nil:
			res.metricField = fm
		}
	}
	return res
}

func newTypeFieldInfo(idx int, sf reflect.StructField) (*typeFieldInfo, *typeMetricField) {
	if !sf.IsExported() && !sf.Anonymous {
		return nil, nil
	}

	tag := sf.Tag.Get(tagName)
	if tag == "-" {
		return nil, nil
	}

	tagMap := ldtags.Parse(tag)
	typ := strings.ToLower(tagMap.Get("type"))
	if typ != "" {
		fnMetric := getMetricValueFuncByType(sf.Type)
		return nil, &typeMetricField{
			Field:        sf,
			Index:        idx,
			Value:        func(v reflect.Value) float64 { return fnMetric(v.Field(idx)) },
			MetricType:   typ,
			MetricAction: tagMap.Get("action"),
			MetricName:   tagMap.Get("metric"),
		}
	}

	p := &typeFieldInfo{
		Field:  sf,
		Index:  idx,
		Name:   tagMap.Get("name"),
		Value:  getValueFuncByType(sf.Type, true),
		Layout: tagMap.Get("layout"),
	}

	if p.Name == "-" || p.Value == nil {
		return nil, nil
	}

	if sf.Anonymous && (sf.Type.Kind() == reflect.Struct || (sf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Struct)) {
		p.Name = ""
	} else if p.Name == "" {
		p.Name = ldstr.ToSnakeCase(sf.Name, '_')
	}

	return p, nil
}

func getMetricValueFuncByType(typ reflect.Type) func(v reflect.Value) float64 {
	v := reflect.Zero(typ)
	switch v.Interface().(type) {
	case time.Duration:
		return func(v reflect.Value) float64 {
			return v.Interface().(time.Duration).Seconds()
		}

	case interface{ Float64() float64 }:
		return func(v reflect.Value) float64 { return v.Interface().(interface{ Float64() float64 }).Float64() }
	case interface{ Float32() float32 }:
		return func(v reflect.Value) float64 {
			return float64(v.Interface().(interface{ Float32() float32 }).Float32())
		}

	case interface{ Int64() int64 }:
		return func(v reflect.Value) float64 { return float64(v.Interface().(interface{ Int64() int64 }).Int64()) }
	case interface{ Int32() int32 }:
		return func(v reflect.Value) float64 { return float64(v.Interface().(interface{ Int32() int32 }).Int32()) }
	case interface{ Int16() int16 }:
		return func(v reflect.Value) float64 { return float64(v.Interface().(interface{ Int16() int16 }).Int16()) }
	case interface{ Int8() int8 }:
		return func(v reflect.Value) float64 { return float64(v.Interface().(interface{ Int8() int8 }).Int8()) }
	case interface{ Int() int }:
		return func(v reflect.Value) float64 { return float64(v.Interface().(interface{ Int() int }).Int()) }

	case interface{ Uint64() uint64 }:
		return func(v reflect.Value) float64 { return float64(v.Interface().(interface{ Uint64() uint64 }).Uint64()) }
	case interface{ Uint32() uint32 }:
		return func(v reflect.Value) float64 { return float64(v.Interface().(interface{ Uint32() uint32 }).Uint32()) }
	case interface{ Uint16() uint16 }:
		return func(v reflect.Value) float64 { return float64(v.Interface().(interface{ Uint16() uint16 }).Uint16()) }
	case interface{ Uint8() uint8 }:
		return func(v reflect.Value) float64 { return float64(v.Interface().(interface{ Uint8() uint8 }).Uint8()) }
	case interface{ Uint() uint }:
		return func(v reflect.Value) float64 { return float64(v.Interface().(interface{ Uint() uint }).Uint()) }
	}

	switch typ.Kind() {
	// case reflect.Bool:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(v reflect.Value) float64 { return float64(v.Int()) }
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return func(v reflect.Value) float64 { return float64(v.Uint()) }
	case reflect.Float32, reflect.Float64:
		return func(v reflect.Value) float64 { return v.Float() }
		// case reflect.String:
		// case reflect.Struct:
		// case reflect.Ptr:
	}

	return func(v reflect.Value) float64 {
		panic(fmt.Errorf("invalid metric field type. type:%s", typ.String()))
	}
}

func getValueFuncByType(typ reflect.Type, needPtrFunc bool) fieldValueFunc {
	switch {
	case typ == typeOfTime:
		return toLabelsByTime
	case typ.Implements(typeOfError):
		return toLabelsByError
	case typ.Implements(typeOfStringer):
		return toLabelsByStringer
	}

	switch typ.Kind() {
	case reflect.Bool:
		return toLabelsByBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return toLabelsByInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return toLabelsByUint
	case reflect.Float32, reflect.Float64:
		return toLabelsByFloat
	case reflect.String:
		return toLabelsByString
	case reflect.Struct:
		return toLabelsByStruct
	case reflect.Ptr:
		if !needPtrFunc {
			return nil
		}
		fn := getValueFuncByType(typ.Elem(), false)
		if fn == nil {
			return nil
		}
		return getPtrValueFunc(fn)
	}
	return nil
}

func getPtrValueFunc(fn fieldValueFunc) fieldValueFunc {
	return func(c *toLabelsContext, v reflect.Value) {
		if !v.IsNil() {
			fn(c, v.Elem())
			return
		}

		typ := v.Type().Elem()
		if typ.Kind() != reflect.Struct {
			c.Set(nilStr)
			return
		}

		v = reflect.Zero(typ)
		fn(c, v)
	}
}

func toLabelsByError(c *toLabelsContext, v reflect.Value) {
	err, _ := v.Interface().(error)
	name := c.GetName()
	c.Labels[fmt.Sprintf("%s_code", name)] = strconv.Itoa(lderr.GetCode(err))
	c.Labels[fmt.Sprintf("%s_message", name)] = lderr.GetMessage(err)
}

func toLabelsByStringer(c *toLabelsContext, v reflect.Value) {
	vv, _ := v.Interface().(stringer)
	if vv == nil {
		c.Set(nilStr)
		return
	}
	c.Set(vv.String())
}

func toLabelsByTime(c *toLabelsContext, v reflect.Value) {
	vv := v.Interface().(time.Time)
	layout := c.Current.Layout
	// log.Printf(" === 1 layout:%s", layout)
	if layout == "" {
		layout = time.RFC3339
	}
	// log.Printf(" === 2 layout:%s", layout)
	c.Set(vv.Format(layout))
}

func toLabelsByString(c *toLabelsContext, v reflect.Value) { c.Set(v.String()) }
func toLabelsByBool(c *toLabelsContext, v reflect.Value)   { c.Set(strconv.FormatBool(v.Bool())) }
func toLabelsByInt(c *toLabelsContext, v reflect.Value)    { c.Set(strconv.FormatInt(v.Int(), 10)) }
func toLabelsByUint(c *toLabelsContext, v reflect.Value)   { c.Set(strconv.FormatUint(v.Uint(), 10)) }
func toLabelsByFloat(c *toLabelsContext, v reflect.Value) {
	c.Set(strconv.FormatFloat(v.Float(), 'g', -1, 64))
}

func toLabelsByStruct(c *toLabelsContext, v reflect.Value) {
	typ := getTypeInfo(v.Type())
	structToLabels(c, v, typ)
}

func structToLabels(c *toLabelsContext, v reflect.Value, typ *typeInfo) {
	curr := c.Current

	if c.Labels == nil {
		c.Labels = make(map[string]string, len(typ.Fields))
	}
	n := len(c.Name)
	for _, f := range typ.Fields {
		if f.Name != "" {
			c.Name = append(c.Name, f.Name)
		}

		c.Current = f
		f.Value(c, v.Field(f.Index))

		c.Name = c.Name[:n]
	}

	c.Current = curr
}
