/*
 * Copyright (C) distroy
 */

package ldlog

import "testing"

func Test_core_buildRecordStr(t *testing.T) {
	tests := []struct {
		name  string
		msg   string
		attrs []Attr
		want  string
	}{
		{
			name:  "panic message",
			msg:   "panic message",
			attrs: []Attr{String("str", "abc"), Reflect("map", map[string]any{"k": "v"})},
			want:  `panic message. str:"abc", map:{"k":"v"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &core{}
			r := &Record{Message: tt.msg}
			r.AddAttrs(tt.attrs...)
			if got := l.buildRecordStr(r); got != tt.want {
				t.Errorf("core.buildRecordStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
