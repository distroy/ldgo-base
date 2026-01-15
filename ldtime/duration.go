/*
 * Copyright (C) distroy
 */

package ldtime

import (
	"encoding/json"
	"time"

	"github.com/distroy/ldgo-base/3rd/yaml"
	"github.com/distroy/ldgo-base/internal/time_"
	"github.com/distroy/ldgo-base/ldconv"
)

var (
	_ json.Marshaler   = (*Duration)(nil)
	_ json.Unmarshaler = (*Duration)(nil)
	_ yaml.Marshaler   = (*Duration)(nil)
	_ yaml.Unmarshaler = (*Duration)(nil)
)

type Duration time.Duration

func (d Duration) Duration() time.Duration { return d.get() }
func (d Duration) get() time.Duration      { return time.Duration(d) }
func (d *Duration) ptr() *time.Duration    { return (*time.Duration)(d) }

func (d Duration) Abs() Duration                { return Duration(d.get().Abs()) }
func (d Duration) Hours() float64               { return d.get().Hours() }
func (d Duration) Microseconds() int64          { return d.get().Microseconds() }
func (d Duration) Milliseconds() int64          { return d.get().Milliseconds() }
func (d Duration) Minutes() float64             { return d.get().Minutes() }
func (d Duration) Nanoseconds() int64           { return d.get().Nanoseconds() }
func (d Duration) Round(m Duration) Duration    { return Duration(d.get().Round(m.get())) }
func (d Duration) Seconds() float64             { return d.get().Seconds() }
func (d Duration) String() string               { return d.get().String() }
func (d Duration) Truncate(m Duration) Duration { return Duration(d.get().Truncate(m.get())) }

func (d Duration) MarshalJSON() ([]byte, error) { return time_.DurationMarshalJson(d.get()) }
func (d *Duration) UnmarshalJSON(b []byte) error {
	dur, err := time_.DurationUnmarshalJson(b)
	if err == nil {
		*d.ptr() = dur
	}
	return err
}

func (d Duration) MarshalYAML() (any, error) { return d.String(), nil }
func (d *Duration) UnmarshalYAML(n *yaml.Node) error {
	// log.Printf(" === UnmarshalYAML: %s", n.Value)
	s := n.Value
	b := ldconv.StrToBytesUnsafe(s)
	if dur, err := time_.ParseDurationByNumber(b); err == nil {
		*d.ptr() = dur
		return nil
	}
	dur, err := time.ParseDuration(s)
	*d.ptr() = dur
	return err
}
