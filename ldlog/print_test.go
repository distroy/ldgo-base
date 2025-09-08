/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"reflect"
	"sort"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func Test_sortedMap(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l0 := []interface{}{
			"abc",
			123,
			int64(234),
			uint(100),
			uint64(200),
			"bcd",
		}
		l1 := make([][2]reflect.Value, 0, len(l0))
		for _, v := range l0 {
			l1 = append(l1, [2]reflect.Value{
				reflect.ValueOf(v), reflect.ValueOf(0),
			})
		}
		sort.Sort(sortedMap(l1))
		// sortedMap(l1).Swap(0, 1)
		l2 := make([]interface{}, 0, len(l0))
		for _, v := range l1 {
			l2 = append(l2, v[0].Interface())
		}

		convey.So(l2, convey.ShouldResemble, []interface{}{
			uint(100),
			123,
			uint64(200),
			int64(234),
			"abc",
			"bcd",
		})
	})
}
