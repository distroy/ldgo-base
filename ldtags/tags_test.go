/*
 * Copyright (C) distroy
 */

package ldtags

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		tag  string
		want Tags
	}{
		{
			tag: `name:size; default:1; meta:n; :x; bool;`,
			want: Tags{
				`name`:    []string{"size"},
				`default`: []string{"1"},
				`meta`:    []string{"n"},
				`bool`:    []string{""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			if got := Parse(tt.tag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
